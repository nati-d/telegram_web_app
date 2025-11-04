import { useEffect, useState } from "react";

declare global {
  interface Window {
    Telegram: any;
  }
}

function App() {
  const [user, setUser] = useState<any>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    const tg = window.Telegram.WebApp;
    tg.ready();

    // Get Telegram user info
    const userData = tg.initDataUnsafe?.user;
    setUser(userData);
  }, []);

  const sendData = async () => {
    if (!user) return;

    const payload = {
      user_id: user.id,
      username: user.username,
      first_name: user.first_name,
      message,
    };

    await fetch("http://localhost:8080/api/submit", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    window.Telegram.WebApp.close(); // optional: closes the WebApp
  };

  return (
    <div style={{ textAlign: "center", padding: "2rem" }}>
      <h2>Hello {user?.first_name || "there"} 👋</h2>
      <p>Send a message to your bot via Go backend!</p>

      <input
        placeholder="Type a message"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        style={{ padding: "8px", marginBottom: "10px", borderRadius: "6px" }}
      />

      <br />

      <button
        onClick={sendData}
        style={{
          padding: "10px 20px",
          borderRadius: "8px",
          border: "none",
          backgroundColor: "#2AABEE",
          color: "white",
          cursor: "pointer",
        }}
      >
        Send to Bot
      </button>
    </div>
  );
}

export default App;
