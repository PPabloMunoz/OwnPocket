# OwnPocket

> Privacy-first self-hosted personal finance manager.

OwnPocket is a lightweight self-hosted finance manager focused on privacy, speed, and ownership. Track expenses, manage budgets, and monitor your finances locally — with your data staying entirely under your control.

---

## Quick Start

### Running with Docker

```bash
docker run -p 8080:8080 -v ./data:/data ppablomunoz/ownpocket
```

### Building from Source

Ensure you have [Go](https://go.dev/) and [Just](https://github.com/casey/just) installed.

1.  **Build the single binary:**
    ```bash
    just build-local
    ```
2.  **Run the app:**
    ```bash
    ./bin/app
    ```
3.  **Access the UI:** [http://localhost:8080](http://localhost:8080)

---

## Features

* Lightweight and fast
* Self-hosted
* Privacy-focused
* Expense tracking
* Budget management
* Account monitoring
* Clean and simple UI
* No ads
* No subscriptions
* No data harvesting

---

## Philosophy

Most finance apps monetize your financial data.

OwnPocket is built on the opposite idea:

* Your data stays with you
* Your server, your rules
* No third-party tracking
* No cloud dependency unless you choose it

---

## Contributing

Contributions, ideas, and feedback are welcome.

Feel free to open an issue or submit a pull request.

---

## License

MIT

