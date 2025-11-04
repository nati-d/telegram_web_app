export {};

declare global {
  interface TelegramWebAppUser {
    id: number;
    first_name: string;
    last_name?: string;
    username?: string;
    language_code?: string;
  }

  interface TelegramThemeParams {
    bg_color?: string;
    text_color?: string;
    hint_color?: string;
    link_color?: string;
    button_color?: string;
    button_text_color?: string;
  }

  interface TelegramWebApp {
    initData: string;
    initDataUnsafe: {
      user?: TelegramWebAppUser;
      query_id?: string;
      auth_date?: string;
      hash?: string;
    };
    colorScheme: "light" | "dark";
    themeParams: TelegramThemeParams;
    sendData: (data: string) => void;
    close: () => void;
    expand: () => void;
  }

  interface Window {
    Telegram: {
      WebApp: TelegramWebApp;
    };
  }
}
