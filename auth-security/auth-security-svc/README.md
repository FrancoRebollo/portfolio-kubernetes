# plantilla_go

## Descripción

Esta plantilla sirve como base para la creación de repositorios para futuras APIs desarrolladas en Go.

La modalidad de trabajo es tomar este repositorio como referencia para el desarrollo de las APIs pero los diseños y especificaciones estarán dados de acuerdo
a lo establecido en cada implementación.

La estructura básica de carpetas es la siguiente:

1. bin: aquí se almacenarán los compilables para la ejecución de pruebas en el entorno local
2. cmd: aquí se encuentran todos los archivos de entrada a la API e incluye los siguiente subdirectorios
    - apis: aquí se encuentran todos los directorios para iniciar el servidor como main.go y el router, además de los handlers junto a sus middlewares particulares
    - config: aquí se encuentra la configuración global del proyecto, como credenciales de conexiones a las bases y datos de la api
    - tmp
    - utils: aquí se encuentran herramientas adicionales para usar a nivel global como funciones de validación de datos, middlewares generales, etc
3. docs: en este directorio se encuentran los archivos que se crean mediante la extensión de swag, para la documentación mediante swaggo
4. internal: posee los directorios que forman la feature o funcionalidad a desarrollar junto con las estructuras de arquitectura hexagonal:
    - feature1:
        - domains
        - ports
        - repository
        - services
    - feature2:
        - domains
        - ports
        - repository
        - services
    - .
    - .
    - .
    - featureN:
        - domains
        - ports
        - repository
        - services

    - storage: posee las configuraciones de conexión y manejo de bases de datos

Luego se encuentran archivos globales que hacen parte de la configuración del repositorio como:
5. .air.toml: este archivo es una copia del archivo .ait.toml.example que posee las indicaciones para levantar en el entorno de desarrollo la aplicación y poder realizar cambios y reiniciar automáticamente el compilable
6. .en: este archivo es copia del .env.example que poseen configuraciones globales del repositorio como credenciales a base de datos
7. .gitignore: es un archivo que ignora directorios y archivos para el versionador gitlab
8. go.mod y go.sum: archivos que son de configuraciones y versiones de librerías que utiliza el repositorio
9. serve.bat: este archivo por lotes permite desplegar localmente la API para el entorno de desarrollo

## Instalaciones y configuraciones

Para instalar todas las dependencias y librerías que utiliza el repositorio se tienen que correr por consola los siguientes comandos:

- Instalar swag:

    ``
    go install github.com/swaggo/swag/cmd/swag@latest
    ``

- Instalar air:

    ``
    go install github.com/air-verse/air@latest
    ``

- Actualizar resto de dependencias:

    ``
    go mod tidy
    ``

## Despliegue en desarrollo

Para poder levantar la API en un entorno local, se tiene que abrir el proyecto en un editor de código, como Visual Studio Code, y ejecutar por consola el comando:

``
    ./serve.bat
``

Si se quiere revisar la documentación realizada en swagger, una vez desplegada la API se tiene que ir a un navegador y colocar en el buscador por ejemplo: http://localhost:8080/docs/index.html

## Despliegue en producción

