import moment from 'moment';

export async function setMomentLocale(locale) {
  switch (locale) {
    case 'az':
      (await import('moment/dist/locale/az'));
      break;
    case 'kk':
      (await import('moment/dist/locale/kk'));
      break;
    case 'ru':
      (await import('moment/dist/locale/ru'));
      break;
    case 'tr':
      (await import('moment/dist/locale/tr'));
      break;
    case 'zh':
      (await import('moment/dist/locale/zh-cn'));
      break;
    case 'en':
      // en is already loaded
      break;
    default:
      console.error(`moment.js locale ${locale} not supported`);
      return;
  }
  moment.locale(locale);
}