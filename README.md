# CobaltBase

**CobaltBase** is an open-source Backend-as-a-Service (BaaS) platform designed for rapid development, flexibility, and scalability. Built using Go with the Chi router, PostgreSQL for database management, and Gorm as the ORM, CobaltBase aims to simplify backend development by providing a customizable and ready-to-use backend service.

## Key Features (In Development)

### **Go-based Backend**

CobaltBase leverages the speed and efficiency of Go to build a highly performant backend. The backend is designed for developers who want to quickly get started with backend development without worrying about the boilerplate.

### **Chi Router**

CobaltBase uses the lightweight and idiomatic Chi router for clean and organized HTTP routing. The router is simple yet powerful, ensuring that you can build and scale your backend efficiently.

### **PostgreSQL**

Harness the power of PostgreSQL, a robust, reliable, and scalable relational database. PostgreSQL is known for its advanced features and flexibility, making it a perfect fit for dynamic and schema-driven applications.

### **Gorm**

CobaltBase integrates Gorm as the ORM layer, allowing you to easily interact with the database. Gorm simplifies CRUD operations, database migrations, and relationship handling, letting you focus on building your application rather than managing database queries.

### **Schema and Rules Management**

Developers using CobaltBase can dynamically create schemas for their data. You define your data structure through an intuitive web portal, specifying fields, data types (e.g., string, integer, email, etc.), and validation rules. This dynamic schema generation enables quick prototyping and adaptability to changing requirements. Additionally, rules can be set to enforce specific constraints and validation logic, ensuring data integrity.

### **Authentication**

CobaltBase supports flexible authentication options. Set up OAuth integration with multiple providers such as Google, GitHub, and Facebook for seamless user authentication, or use traditional email-password authentication. The platform offers a variety of built-in authentication mechanisms, allowing you to focus on building features rather than handling the complexities of authentication.

### **File Uploads**

The platform includes file upload functionality that integrates with your schemas. CobaltBase manages the file handling for you, storing the files, generating URLs, and associating them with your data models. Whether you need to handle images, documents, or any other file types, CobaltBase makes file management simple and efficient.

### **Realtime APIs**

CobaltBase supports real-time communication through both Server-Sent Events (SSE) and WebSockets. Developers can build real-time applications like chat apps, notifications, or live updates without needing to worry about the infrastructure. The system automatically handles subscriptions and broadcasts events in real-time, empowering developers to create engaging and interactive applications.

### **RESTful API Generation**

For every schema created, CobaltBase automatically generates RESTful APIs. You can interact with your data using standard RESTful methods (GET, POST, PUT, DELETE), and the API enforces the rules and validation constraints set up in your schema. This auto-generation of APIs saves time and ensures consistency across your application.

### **JavaScript SDK**

CobaltBase offers a JavaScript SDK to help you seamlessly integrate your backend with front-end frameworks like React, Vue, Svelte, and more. The SDK provides convenient methods to interact with your backend, handle authentication, manage files, and subscribe to real-time events, making it easier to build full-stack applications.

### **Web Console**

CobaltBase includes an intuitive web console built with SvelteKit, allowing you to administer and manage your backend directly from the browser. The console provides interfaces for creating schemas, managing authentication, configuring rules, monitoring real-time events, and interacting with your data.

### **Extensibility**

While CobaltBase offers a range of built-in features like CRUD operations, authentication, and real-time APIs, extensibility is intentionally limited. Developers can use CobaltBase as a library, adding their own custom logic through Go's Chi endpoints, middleware, and Gorm queries. This approach allows for flexibility in extending functionality without overcomplicating the platform, ensuring that developers can tailor CobaltBase to their needs while still benefiting from its robust core features.

## Status

CobaltBase is currently in the **development phase**, with many of these features actively being worked on. The platform is evolving to meet the needs of developers seeking a robust and flexible BaaS solution. Stay tuned for more updates as the development progresses.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
