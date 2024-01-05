# simpleapp

## How test on the local machine

1. build+start (detached) via Makefile command

    ```shell
    make docker_compose_up
    ```

2. test application via curl command

    ```shell
    $ curl localhost:8080
    ```

    curl output

    ```html
    <!DOCTYPE html>
    <html>
            <head>
                    <title>Simple App</title>
            </head>
            <body>
                    <h1>Simple App</h1>
                    <p>The env <b>$MESSAGE</b> of the backend system is: <b>Hello from backend</b></p>
            </body>
    </html>
    ```
3. stop local stage

    ```
    make docker_compose_down
    ```

4. cleanup (alternative to point 3)

    ```
    make clean
    ```
