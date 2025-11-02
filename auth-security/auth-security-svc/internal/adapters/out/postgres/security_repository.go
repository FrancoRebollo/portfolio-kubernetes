package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/FrancoRebollo/auth-security-svc/internal/domain"
	"github.com/FrancoRebollo/auth-security-svc/internal/platform/utils"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type SecurityRepository struct {
	dbPost *PostgresDB
}

func (v SecurityRepository) CheckAPI2FA(ctx context.Context, idPersona int, apiKey string, canalDigital string) (*string, error) {
	var reqApiKey string
	var reqUser string
	var seed2FAPointer *string
	var username string
	var seed2FAString string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT req_2fa from sec.api_key where api_key = $1 `

	err = tx.QueryRowContext(ctx, query, apiKey).Scan(&reqApiKey)

	if err != nil {
		return nil, err
	}

	query = `SELECT req_2fa,login_name from sec.canal_digital_persona where id_persona = $1 and tipo_canal_digital = $2 `

	err = tx.QueryRowContext(ctx, query, idPersona, canalDigital).Scan(&reqUser, &username)

	if err != nil {
		return nil, err
	}

	if reqUser == "N" && reqApiKey == "N" {
		return nil, nil
	}

	query = `SELECT coalesce("2fa_seed",'NO TIENE') from sec.token where api_key = $1 and id_canal_digital_persona =
		(select id_canal_digital_persona from sec.canal_digital_persona where id_persona = $2 and tipo_canal_digital = $3)`

	err = tx.QueryRowContext(ctx, query, apiKey, idPersona, canalDigital).Scan(&seed2FAString)

	if err != nil {
		return nil, err
	}

	if seed2FAString == "NO TIENE" {

		seed2FA, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Thinksoft-autenticacion",
			AccountName: username,
		})

		if err != nil {
			return nil, err
		}

		update := `update sec.token set "2fa_seed" = $1 where api_key = $1 and id_canal_digital_persona =
		(select id_canal_digital_persona from sec.canal_digital_persona where id_persona = $2 and tipo_canal_digital = $3)`

		_, err = tx.ExecContext(ctx, update, seed2FA.Secret(), idPersona, canalDigital)

		if err != nil {
			return nil, err
		}

		seed2FAString = seed2FA.Secret()

		seed2FAPointer = &seed2FAString
	}

	if err = tx.Commit(); err != nil {
		return seed2FAPointer, err
	}

	seed2FAPointer = &seed2FAString

	return seed2FAPointer, nil
}

func (v SecurityRepository) checkRevokes(ctx context.Context, credentials domain.Credentials) error {

	query := `SELECT id_persona FROM sec.persona WHERE id_persona = $1 and acceso_revocado = 'S'`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("persona revocada")
	}
	rows.Close()

	query = `SELECT tipo_canal_digital FROM sec.tipo_canal_digital_df WHERE tipo_canal_digital = $1 and acceso_revocado = 'S'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("canal digital revocado")
	}
	rows.Close()

	query = `SELECT id_canal_digital_persona FROM sec.canal_digital_persona WHERE tipo_canal_digital = $1
		and id_persona = $2 and canal_validado = 'N' and acceso_revocado = 'S'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital, credentials.IdPersona)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("acceso revocado persona - canal digital")
	}
	rows.Close()

	query = `SELECT api_key FROM sec.api_key WHERE api_key = $1 and fecha_fin_vigencia < current_date`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.ApiKey)

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("api key expirada")
	}
	rows.Close()

	return nil
}

func (v SecurityRepository) checkCredentials(ctx context.Context, credentials domain.Credentials) error {

	query := `SELECT id_persona FROM sec.persona WHERE id_persona = $1`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("persona no encontrada")
	}
	rows.Close()

	query = `SELECT tipo_canal_digital FROM sec.tipo_canal_digital_df WHERE tipo_canal_digital = $1`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("canal digital invalido")
	}
	rows.Close()

	query = `SELECT id_canal_digital_persona FROM sec.canal_digital_persona WHERE tipo_canal_digital = $1
		and id_persona = $2 and canal_validado = 'N'`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.CanalDigital, credentials.IdPersona)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("canal digital no validado")
	}
	rows.Close()

	query = `SELECT api_key FROM sec.api_key WHERE api_key = $1`

	rows, err = v.dbPost.GetDB().QueryContext(ctx, query, credentials.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("api key desconocida")
	}
	rows.Close()

	return nil
}

// CreateUser implements ports.SecurityRepository.
func (s *SecurityRepository) CreateUser(ctx context.Context, reqAltaUser domain.UserCreated) (*domain.UserCreated, error) {

	var message string
	var idPersona int
	tx, err := s.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `SELECT TIPO_CANAL_DIGITAL FROM sec.TIPO_CANAL_DIGITAL_DF WHERE TIPO_CANAL_DIGITAL = $1`

	rows, err := tx.QueryContext(ctx, query, reqAltaUser.CanalDigital)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("canal digital invalido")
	}
	rows.Close()

	query = `SELECT id_persona FROM sec.persona WHERE id_persona = $1`

	err = tx.QueryRowContext(ctx,
		`SELECT id_persona FROM sec.persona WHERE id_persona = $1`,
		reqAltaUser.IdPersona, // pass as int, no fmt.Sprint
	).Scan(&idPersona)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = tx.QueryRowContext(ctx,
				`INSERT INTO sec.persona (last_location) VALUES ($1) RETURNING id_persona`,
				"1",
			).Scan(&idPersona)
			if err != nil {
				return nil, err
			}

			reqAltaUser.IdPersona = idPersona
		} else {
			return nil, err
		}
	}

	rows.Close()

	query = `SELECT id_canal_digital_persona FROM sec.canal_digital_persona WHERE id_persona = $1 and tipo_canal_digital = $2`

	rows, err = tx.QueryContext(ctx, query, fmt.Sprint(reqAltaUser.IdPersona), reqAltaUser.CanalDigital)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqAltaUser.Password), bcrypt.DefaultCost)

		insert := `INSERT INTO sec.CANAL_DIGITAL_PERSONA 
		(ID_PERSONA,TIPO_CANAL_DIGITAL,PASSWORD_ACCESO_HASH,mail_persona,telefono_persona,LOGIN_NAME) VALUES ($1,$2,$3,$4,$5,$6)`

		_, err = tx.ExecContext(ctx, insert, fmt.Sprint(reqAltaUser.IdPersona), reqAltaUser.CanalDigital, hashedPassword, reqAltaUser.MailPersona,
			reqAltaUser.TePersona, reqAltaUser.LoginName)

		if err != nil {
			return nil, err
		}

		message += " en su canal digital"
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.UserCreated{
		IdPersona:    reqAltaUser.IdPersona,
		CanalDigital: reqAltaUser.CanalDigital,
		LoginName:    reqAltaUser.LoginName,
		Password:     reqAltaUser.Password,
		MailPersona:  reqAltaUser.MailPersona,
		TePersona:    reqAltaUser.TePersona,
	}, nil
}

func NewSecurityRepository(dbPost *PostgresDB) *SecurityRepository {
	return &SecurityRepository{
		dbPost: dbPost,
	}
}

func (v SecurityRepository) checkSuperUser(ctx context.Context, apiKey string) (error, bool) {
	var isSuperUser string

	query := `select is_super_user from sec.api_key where api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&isSuperUser)

	if err != nil {
		return err, false
	}

	if isSuperUser == "N" {
		return nil, false
	}

	return nil, true
}

func (v SecurityRepository) CrearCanalDigital(ctx context.Context, crearCanalDigital domain.CanalDigital, apiKey string) error {

	err, isSuperUser := v.checkSuperUser(ctx, apiKey)

	if err != nil {
		return fmt.Errorf("Unknown api-key")
	}

	if !isSuperUser {
		return fmt.Errorf("no posee los permisos necesarios para esta operacion")
	}

	insert := `insert into sec.tipo_canal_digital_df (tipo_canal_digital) values ($1)`

	_, err = v.dbPost.GetDB().ExecContext(ctx, insert, crearCanalDigital.CanalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) AccessPerson(ctx context.Context, accessPerson domain.AccessPerson, apikey string) error {

	update := `update sec.persona set acceso_revocado = $1 where id_persona = $2`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, accessPerson.Revoke, accessPerson.IdPersona)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) AccessCanalDigital(ctx context.Context, accessCanaldigital domain.AccessCanalDigital, apikey string) error {
	update := `update sec.tipo_canal_digital_df set acceso_revocado = $1 where tipo_canal_digital = $2`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, accessCanaldigital.Revoke, accessCanaldigital.CanalDigital)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) AccessPersonMethodAuth(ctx context.Context, accesPersonMethodAuth domain.AccessPersonMethodAuth, apikey string) error {
	update := `update sec.canal_digital_persona set acceso_revocado = $1 where tipo_canal_digital = $2 
		and id_persona = $3`

	_, err := v.dbPost.GetDB().ExecContext(ctx, update, accesPersonMethodAuth.Revoke, accesPersonMethodAuth.MethodAuth, accesPersonMethodAuth.IdPersona)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) AccessApiKey(ctx context.Context, accessApiKey domain.AccessApiKey, apikey string) error {

	if accessApiKey.Revoke == "N" {
		layout := accessApiKey.FechaVigencia
		t, err := time.ParseInLocation(layout, "12/10/2025 00:00:00", time.FixedZone("ART", -3*3600))
		if err != nil {
			return err
		}

		t = t.UTC()
		update := `update sec.api_key set fecha_vigencia = $1 where api_key = $2`

		_, err = v.dbPost.GetDB().ExecContext(ctx, update, t, accessApiKey.ApiKey)

		if err != nil {
			return err
		}
	} else {
		update := `update sec.api_key set fecha_vigencia = $1 where api_key = $2`

		_, err := v.dbPost.GetDB().ExecContext(ctx, update, "CURRENT_DATE - INTERVAL '1 DAY'", accessApiKey.ApiKey)

		if err != nil {
			return err
		}
	}

	return nil
}

func (v SecurityRepository) LoginValidations(ctx context.Context, reqLogin domain.Login) (int, *string, error) {

	var idPersona int
	var hashedPassword string

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return 0, nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `select id_persona,password_acceso_hash from sec.canal_digital_persona where tipo_canal_digital = $1 
		and login_name = $2`

	err = tx.QueryRowContext(ctx, query, reqLogin.CanalDigital, reqLogin.Username).Scan(&idPersona, &hashedPassword)

	if idPersona == 0 {
		return 0, nil, fmt.Errorf("usuario o canal digital incorrecto")
	}

	if err != nil {
		return idPersona, nil, err
	}

	if err = utils.ComparePasswordHash(hashedPassword, reqLogin.Password); err != nil {
		return idPersona, nil, fmt.Errorf("contraseña incorrecta")
	}

	credentials := domain.Credentials{
		IdPersona:    idPersona,
		ApiKey:       reqLogin.ApiKey,
		CanalDigital: reqLogin.CanalDigital,
	}

	if err := v.checkCredentials(ctx, credentials); err != nil {
		return idPersona, nil, err
	}

	if err := v.checkRevokes(ctx, credentials); err != nil {
		return idPersona, nil, err
	}

	seed2FA, err := v.CheckAPI2FA(ctx, idPersona, reqLogin.ApiKey, reqLogin.CanalDigital)

	if err != nil {
		return idPersona, nil, err
	}

	if err = tx.Commit(); err != nil {
		return idPersona, seed2FA, err
	}

	return idPersona, seed2FA, nil
}

func (v SecurityRepository) UpsertAccessToken(ctx context.Context, requestUpsert *domain.UpsertAccessToken) error {
	var idCanalDigitalPersona int

	expAccessToken, err := utils.GetTokenExpiration(requestUpsert.AccessToken, "ACCESS")

	if err != nil {
		return err
	}

	expRefreshToken, err := utils.GetTokenExpiration(requestUpsert.RefreshToken, "REFRESH")

	if err != nil {
		return err
	}

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT ID_CANAL_DIGITAL_PERSONA FROM sec.CANAL_DIGITAL_PERSONA
		WHERE ID_PERSONA = $1 AND TIPO_CANAL_DIGITAL = $2`

	err = tx.QueryRowContext(ctx, query, requestUpsert.IdPersona, requestUpsert.CanalDigital).Scan(&idCanalDigitalPersona)

	if err != nil {
		return err
	}

	query = `SELECT id_token FROM sec.TOKEN WHERE ID_CANAL_DIGITAL_PERSONA = $1 AND API_KEY = $2`

	rows, err := tx.QueryContext(ctx, query, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {

		insert := `INSERT INTO sec.token 
		(api_key,id_canal_digital_persona,access_token,fecha_exp_access_token,refresh_token,fecha_Exp_refresh_token) VALUES ($1,$2,$3,$4,$5,$6)`

		_, err = tx.ExecContext(ctx, insert, requestUpsert.ApiKey, idCanalDigitalPersona, requestUpsert.AccessToken,
			expAccessToken, requestUpsert.RefreshToken, expRefreshToken)

		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return nil
	}
	rows.Close()

	insert := `
    INSERT INTO sec.hist_token (id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado)
    SELECT id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado
    FROM sec.token
    WHERE id_canal_digital_persona = $1
	and api_key = $2`

	_, err = tx.ExecContext(ctx, insert, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	update := `update sec.token	set access_token = $1, fecha_creacion_token = $2, fecha_exp_access_token = $3, refresh_token = $4
		,fecha_exp_refresh_token = $5 
		where id_canal_digital_persona = $6 
		and api_key = $7`

	_, err = tx.ExecContext(ctx, update, requestUpsert.AccessToken, time.Now(), expAccessToken, requestUpsert.RefreshToken,
		expRefreshToken, idCanalDigitalPersona, requestUpsert.ApiKey)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) GetAccessTokenDuration(ctx context.Context, ApiKey string) (int, error) {

	var ctdHoras int

	query := `SELECT ctd_hs_access_token_valido FROM sec.api_key
		WHERE api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, ApiKey).Scan(&ctdHoras)

	if err != nil {
		return 0, err
	}

	return ctdHoras * 60, nil
}

func (v SecurityRepository) CheckTokenCreation(ctx context.Context, credentials domain.Credentials) error {

	if err := v.checkCredentials(ctx, credentials); err != nil {
		return err
	}

	if err := v.checkRevokes(ctx, credentials); err != nil {
		return err
	}

	query := `SELECT id_token,fecha_exp_refresh_token FROM sec.TOKEN
		WHERE ID_CANAL_DIGITAL_PERSONA = (SELECT ID_CANAL_DIGITAL_PERSONA FROM sec.CANAL_DIGITAL_PERSONA
											WHERE ID_PERSONA = $1 AND TIPO_CANAL_DIGITAL = $2)
		and api_key = $3`

	rows, err := v.dbPost.GetDB().QueryContext(ctx, query, credentials.IdPersona, credentials.CanalDigital, credentials.ApiKey)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf("loguee por primera vez")
	}

	rows.Close()

	return nil
}

func (v SecurityRepository) CheckLastRefreshToken(ctx context.Context, token string, credentials domain.Credentials) error {

	var idToken int

	query := `SELECT id_token FROM sec.token
		WHERE api_key = $1 
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from sec.canal_digital_persona
											where id_persona = $2
											and tipo_canal_digital = $3)
		and refresh_token = $4`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.ApiKey, credentials.IdPersona, credentials.CanalDigital,
		token).Scan(&idToken)

	if err != nil {
		return fmt.Errorf("token de refresco desconocido")
	}

	return nil
}

func (v SecurityRepository) CheckLastAccessToken(ctx context.Context, token string, credentials domain.Credentials) error {

	var idToken int

	query := `SELECT id_token FROM sec.token
		WHERE api_key = $1 
		AND id_canal_digital_persona = (select id_canal_digital_persona 
											from sec.canal_digital_persona
											where id_persona = $2
											and tipo_canal_digital = $3)
		and access_token = $4`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.ApiKey, credentials.IdPersona, credentials.CanalDigital,
		token).Scan(&idToken)

	if err != nil {
		return fmt.Errorf("el token de acceso no coincide con el ultimo registrado")
	}

	return nil
}

func (v SecurityRepository) updateAccessToken(ctx context.Context, credentialsToken domain.CredentialsToken, idCanalDigitalPersona int) error {

	var idToken int
	accesTokenDuration, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION"))

	if err != nil {
		return fmt.Errorf("error calculando duracion de refresh token")
	}

	accessExpiresAt := time.Now().Add(time.Minute * time.Duration(accesTokenDuration))

	query := `SELECT id_token FROM sec.TOKEN	WHERE ID_CANAL_DIGITAL_PERSONA = $1 AND API_KEY = $2`

	_ = v.dbPost.GetDB().QueryRowContext(ctx, query, idCanalDigitalPersona, credentialsToken.ApiKey).Scan(&idToken)

	if idToken == 0 {
		return fmt.Errorf("registro no encontrado en token - loguee por primera vez")
	}

	insert := `
    INSERT INTO sec.hist_token (id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado)
    SELECT id_token, api_key, id_canal_digital_persona, access_token, fecha_creacion_token, fecha_exp_access_token, refresh_token, fecha_exp_refresh_token, acceso_revocado
    FROM sec.token
    WHERE id_token = $1`

	_, err = v.dbPost.GetDB().ExecContext(ctx, insert, idToken)

	if err != nil {
		return err
	}

	update := `update sec.token set access_token = $1, fecha_exp_access_token = $2 
		where id_token = $3`

	_, err = v.dbPost.GetDB().ExecContext(ctx, update, credentialsToken.AccessToken, accessExpiresAt, idToken)

	if err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) PersistToken(ctx context.Context, credentials domain.CredentialsToken) error {
	query := `SELECT ID_CANAL_DIGITAL_PERSONA FROM sec.CANAL_DIGITAL_PERSONA WHERE ID_PERSONA = $1 
		AND TIPO_CANAL_DIGITAL = $2`

	var idCanalDigitalPersona int

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, credentials.IdPersona, credentials.CanalDigital).Scan(&idCanalDigitalPersona)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no se encontró el canal digital en relacion a la persona")
		}
		return err
	}

	if err = v.updateAccessToken(ctx, credentials, idCanalDigitalPersona); err != nil {
		return err
	}

	return nil
}

func (v SecurityRepository) CheckApiKeyExpirada(ctx context.Context, apiKey string) (bool, error) {
	var fecha time.Time
	var apiKeyRevocada string

	query := `SELECT COALESCE(fecha_fin_vigencia, CURRENT_DATE + INTERVAL '10 days') FROM sec.api_key WHERE api_key = $1`

	err := v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&fecha)

	if err != nil {
		return false, err
	}

	if fecha.Before(time.Now()) {
		return true, nil
	}

	query = `SELECT estado FROM sec.api_key WHERE api_key = $1`

	err = v.dbPost.GetDB().QueryRowContext(ctx, query, apiKey).Scan(&apiKeyRevocada)

	if err != nil {
		return false, err
	}

	if apiKeyRevocada == "INACTIVO" {
		return true, nil
	}

	return false, nil
}

func (v SecurityRepository) CambioPasswordByLogin(ctx context.Context, loginName string, newPassword string) error {

	tx, err := v.dbPost.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT mail_persona FROM sec.CANAL_DIGITAL_PERSONA WHERE login_name = $1`

	var mailPersona string
	err = v.dbPost.GetDB().QueryRowContext(ctx, query, loginName).Scan(&mailPersona)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no se encontró el canal digital en relacion a la persona")
		}
		return err
	}

	if mailPersona == "" {
		return fmt.Errorf("no fue posible recuperar el correo para la recuperacion")
	}

	body := fmt.Sprintf("Hola, tu nueva contraseña es: %s", newPassword)

	if err = utils.SendEmail(mailPersona, "Cambio de contraseña", body); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	update := `update sec.canal_digital_persona set password_acceso_hash = $1 where login_name = $2`

	_, err = tx.ExecContext(ctx, update, hashedPassword, loginName)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
