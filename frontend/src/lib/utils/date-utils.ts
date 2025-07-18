import { APP_CONFIG } from "$lib/constants/config";

export function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString(APP_CONFIG.locales.language, APP_CONFIG.locales.defaultDateFormat);
}
