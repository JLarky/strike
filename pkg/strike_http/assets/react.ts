export declare function useState<T>(
  initialValue: () => T
): [T, (newValue: T | ((oldValue: T) => T)) => void];
export declare function useState<T>(
  initialValue: T
): [T, (newValue: T) => void];

declare module "react/jsx-runtime" {
  export function jsx(
    type: any,
    props?: any,
    key?: string | number | null,
    isStaticChildren?: boolean
  ): any;
  export function jsxs(
    type: any,
    props?: any,
    key?: string | number | null,
    isStaticChildren?: boolean
  ): any;
}
