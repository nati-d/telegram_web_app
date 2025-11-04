import { useEffect, useState } from "react";

declare global {
  interface Window {
    Telegram: any;
  }
}

function App() {
  const [themeParams, setThemeParams] = useState<any>({});
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const tg = window.Telegram.WebApp;
    tg.ready(); // let Telegram know the app is ready

    // Set initial theme parameters
    setThemeParams(tg.themeParams);
    setIsDark(tg.colorScheme === "dark");

    // Listen for theme change events
    tg.onEvent("themeChanged", () => {
      setThemeParams(tg.themeParams);
      setIsDark(tg.colorScheme === "dark");
    });

    return () => tg.offEvent("themeChanged");
  }, []);

  return (
    <div
      style={{
        backgroundColor: isDark ? themeParams.bg_color : "#fff",
        color: isDark ? themeParams.text_color : "#000",
        minHeight: "100vh",
        padding: "2rem",
        textAlign: "center",
        transition: "all 0.3s ease",
      }}
    >
      <h1>Telegram Theme Demo 🧠</h1>
      <p>
        This app automatically matches your Telegram’s{" "}
        {isDark ? "dark" : "light"} mode.
      </p>
      <button
        style={{
          backgroundColor: themeParams.button_color || "#0088cc",
          color: themeParams.button_text_color || "#fff",
          border: "none",
          padding: "10px 20px",
          borderRadius: "8px",
          cursor: "pointer",
        }}
      >
        Telegram Style Button
      </button>
    </div>
  );
}

export default App;
