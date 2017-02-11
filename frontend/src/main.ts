import "./polyfills.browser";

import { platformBrowser } from "@angular/platform-browser";
import { AppModule } from "./app/app.module";

export const platformRef = platformBrowser();

export function main() {
  return platformRef.bootstrapModule(AppModule)
    .catch(err => console.error(err));
}

// support async tag or hmr
switch (document.readyState) {
  case "interactive":
  case "complete":
    main();
    break;
  case "loading":
  default:
    document.addEventListener("DOMContentLoaded", () => main());
}
