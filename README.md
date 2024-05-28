# PassFort

PassFort is a secure password and secret management application designed to store and manage sensitive information such as passwords and text secrets. The application ensures that secrets are stored securely and only accessible by authorized users.

## Features

- User registration with email and master password
- Secure storage of secrets in collections
- Supports two types of secrets:
  - Password secrets
  - Text secrets
- RESTful API for managing secrets and collections
- User authentication and authorization

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/8thgencore/passfort.git
   cd passfort
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Set up environment variables:

   Create a `.env` file in the project root and add the necessary environment variables.

   ```env
    APP_NAME=passfort

    DB_CONNECTION=postgres
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=passfort
    DB_USER=user
    DB_PASSWORD=password

    REDIS_PASSWORD=password

    TOKEN_SIGNING_KEY=eab9adf86028e9409c785431114c9426
   ```

4. Run the application:

   ```sh
   task dev
   ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING](CONTRIBUTING.md) guidelines before submitting a pull request.

## Contact

For any questions or support, please contact us at support@passfort.com.

---

This README provides a basic overview of the PassFort project, including installation instructions, key features, API endpoints, and other essential information. For more detailed documentation, please refer to the project files and additional documentation resources.
