# Waguri Token Joiner

A high-performance, multi-threaded Discord token joiner using uTLS and HTTP/2 for maximum stealth and speed.

---

## 🚀 Features

- ✅ Token format support: `token` or `email:pass:token`
- ✅ Configurable thread count and invite link via `config.json`
- ✅ Supports both proxy and proxyless modes
- ✅ Modular headers, TLS, and logging
- ✅ Auto-handles gzip/deflate responses
- ✅ Fully customizable for mass-join tools

---

## 📦 Requirements

- [Go (Golang)](https://golang.org/dl/) 1.18+
- Git
- A terminal or shell (PowerShell, CMD, Bash)
- Proxy list (optional) in format: `ip:port:user:pass`

---

## 🔧 Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mascular/Discord-Token-Joiner.git
   cd waguri-joiner
   ```

2. **Install dependencies**

   ```bash
   go get ./...
   ```

3. **Edit `config.json`**

   ```json
   {
     "proxy": "", Empty For Proxyless
     "threads": 5,
     "invite": "your_invite_code" 
   }
   ```

4. **Add your tokens**
   Add them to `tokens.txt`, either as:

   * Just token:

     ```
     mfa.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
     ```
   * Or with email/pass:

     ```
     email@example.com:password:mfa.xxxxxxxxxxxxxxxxxxxxxxxxxxxx
     ```

---

## ▶️ Run the Joiner

```bash
go run main.go
```

---

## 🌐 Support & Community

Need help or want updates?
Join the official Waguri support Discord:

🔗 [**discord.gg/waguri-san**](https://discord.gg/waguri-san)

---
