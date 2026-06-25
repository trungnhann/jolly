import enEN from "../locales/en_EN.json";

// Type definitions to infer nested keys
type NestedKeyOf<ObjectType extends object> = {
  [Key in keyof ObjectType & (string | number)]: ObjectType[Key] extends object
    ? `${Key}.${NestedKeyOf<ObjectType[Key]>}`
    : `${Key}`;
}[keyof ObjectType & (string | number)];

export type Translations = typeof enEN;
export type TranslationKey = NestedKeyOf<Translations>;

export function getTranslation(key: TranslationKey): string {
  const keys = key.split(".");
  let current: any = enEN;
  for (const k of keys) {
    if (current[k] === undefined) return key;
    current = current[k];
  }
  return current;
}

export function useTranslation() {
  return {
    t: getTranslation,
  };
}
