export function useState<T>(initialValue: () => T): [T, (newValue: T) => void];
export function useState<T>(initialValue: T): [T, (newValue: T) => void];
