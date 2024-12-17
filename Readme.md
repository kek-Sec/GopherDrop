# 🛠️ **GopherDrop** – Secure One-Time Secret Sharing 🏁

GopherDrop is a secure, self-hostable REST API and UI for sharing encrypted one-time secrets and files, inspired by Bitwarden's Send feature. Built with **Go**, **Vue.js**, and **Vuetify**, GopherDrop is simple to deploy, easy to use, and designed with security in mind.

![GopherDrop Banner](ui/src/assets/Images/banner.png)

---

## 🌟 **Features**

- **Send Text or Files**: Share sensitive information securely.
- **Password Protection**: Encrypt your secrets with a password.
- **One-Time Retrieval**: Automatically delete secrets after they're accessed once.
- **Expiration Settings**: Define how long a secret remains available.
- **Responsive UI**: Built with Vue.js and Vuetify for a clean, modern look.
- **Dockerized Deployment**: Simple setup with Docker and Docker Compose.

---

## 🚀 **Getting Started**

### **Prerequisites**

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)

### **Clone the Repository**

```bash
git clone https://github.com/kek-Sec/gopherdrop.git
cd gopherdrop
```

---

## 🛠️ **Build and Run**

### **Development Setup**

To run GopherDrop locally for development:

```bash
make build      # Build the Docker images
make up         # Start the backend, frontend, and database services
```

The application will be available at:

- **UI**: `http://localhost:8081`
- **API**: `http://localhost:8080`

### **Stopping Services**

```bash
make down
```

### **Running Tests**

```bash
make test
```

---

## ⚙️ **Configuration**

GopherDrop uses environment variables for configuration. These can be set in the `docker-compose.yml` file.

### **Environment Variables**

| Variable         | Description                     | Default Value                        |
|------------------|---------------------------------|--------------------------------------|
| `DB_HOST`        | Database host                   | `db`                                |
| `DB_USER`        | Database username               | `user`                              |
| `DB_PASSWORD`    | Database password               | `pass`                              |
| `DB_NAME`        | Database name                   | `gopherdropdb`                      |
| `SECRET_KEY`     | Secret key for encryption       | `supersecretkeysupersecretkey32`    |
| `LISTEN_ADDR`    | API listen address              | `:8080`                             |
| `STORAGE_PATH`   | Path for storing uploaded files | `/app/storage`                      |
| `MAX_FILE_SIZE`  | Maximum file size in bytes      | `10485760` (10 MB)                  |

---

## 🖥️ **Endpoints**

### **API Endpoints**

| Method | Endpoint           | Description                              |
|--------|--------------------|------------------------------------------|
| `POST` | `/send`            | Create a new send (text or file)         |
| `GET`  | `/send/:id`        | Retrieve a send by its hash              |
| `GET`  | `/send/:id/check`  | Check if a send requires a password      |

### **Example: Create a Send**

```bash
curl -X POST http://localhost:8080/send \
  -F "type=text" \
  -F "data=This is a secret message" \
  -F "password=mysecurepassword"
```

---

## 🌐 **UI**

The UI provides a simple way to create and view sends.

### **Create a Send**

1. Go to `http://localhost:8081`
2. Enter your secret (text or file).
3. Set an optional password and expiration.
4. Click **Create** and share the generated link.

### **View a Send**

1. Open the provided link.
2. If password-protected, enter the password.
3. View or copy the secret. Download files directly from the UI.

---

## 🐳 **Docker Deployment**

To deploy GopherDrop using Docker Compose:

```bash
make build
make up
```

---

## 🔒 **Security Best Practices**

- **Encryption**: All data is encrypted with AES-256 before storage.
- **Passwords**: Optional password protection for additional security.
- **One-Time Use**: Option to delete sends after a single access.
- **Expiration**: Sends auto-delete after the set expiration time.

---

## 🤝 **Contributing**

We welcome contributions from the community! To contribute:

1. Fork the repository.
2. Create a new branch: `git checkout -b my-feature-branch`
3. Make your changes and add tests.
4. Submit a pull request.

### **Contribution Guidelines**

- Follow the project’s coding standards.
- Keep your code clean and well-documented.
- Ensure tests pass before submitting your PR.

---

## 📝 **License**

GopherDrop is licensed under the [MIT License](LICENSE).

---

## 💬 **Community and Support**

- **Issues**: [GitHub Issues](https://github.com/kek-Sec/gopherdrop/issues)
- **Discussions**: Join the discussion in our [GitHub Discussions](https://github.com/kek-Sec/gopherdrop/discussions)
