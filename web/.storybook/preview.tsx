import type { Decorator, Preview } from "@storybook/react-vite";
import "../src/app/styles/index.css";

/** Toolbar de tema: cada story se ve en dark (default de marca) y light (respiro). */
const withTheme: Decorator = (Story, context) => {
  const theme = (context.globals.theme as string) ?? "dark";
  document.documentElement.setAttribute("data-theme", theme);
  return (
    <div className="bg-background p-6 text-foreground" style={{ minHeight: "100px" }}>
      <Story />
    </div>
  );
};

const preview: Preview = {
  globalTypes: {
    theme: {
      description: "Tema PRENTER",
      toolbar: {
        title: "Tema",
        items: [
          { value: "dark", title: "Dark (marca)" },
          { value: "light", title: "Light (respiro)" },
        ],
        dynamicTitle: true,
      },
    },
  },
  initialGlobals: { theme: "dark" },
  decorators: [withTheme],
};

export default preview;
