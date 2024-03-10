## Running the Application

### GoLang Backend

To run the GoLang backend, execute the following command:

```bash
go run cmd/main.go
```

This will start the backend server.

## PostgreSQL Database
Make sure to configure your PostgreSQL database by providing the necessary environment variables in a .env file. Example .env file:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name
```

## Vue.js Frontend
The Vue.js frontend is located in the /web directory. During development, you can run the Vue.js development server using:

```bash
cd web
npm install
npm run serve
```

This will start the Vue.js development server. The output directory for production-ready files is /static.
