function snakeToCamel(str: string): string {
  return str.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
}

export function keysToCamel<T>(value: unknown): T {
  if (Array.isArray(value)) {
    return value.map(keysToCamel) as unknown as T;
  }
  if (value !== null && typeof value === 'object') {
    const result: Record<string, unknown> = {};
    for (const key of Object.keys(value as object)) {
      result[snakeToCamel(key)] = keysToCamel((value as Record<string, unknown>)[key]);
    }
    return result as T;
  }
  return value as T;
}
