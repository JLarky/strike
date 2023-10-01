import type { RemotePromise } from "./rsc";

export declare function Router(): JSX.Element;

export type ActionData = {
  actionId: string;
  data: any;
  remotePromise: RemotePromise;
};

export type RouterState = {
  key: string;
  href: string;
  isInitial: boolean;
  actionData?: ActionData;
};

export declare function createRouterState(href: string): RouterState;
export declare function changeRouterState(
  href: string,
  key: string
): RouterState;
export declare function changeRouterStateForAction(
  href: string,
  key: string,
  actionData: ActionData
): RouterState;

export declare function addNavigation(
  setRouter: (
    router: RouterState | ((router: RouterState) => RouterState)
  ) => void
): void;
export declare function navigate(href: string): void;
export declare function submitForm(actionData: ActionData): void;

declare global {
  interface Window {
    __rsc: any;
    __rscNav: (pathname: string) => void;
    __rscAction: <T>(
      formAction: string,
      formData: FormData | undefined
    ) => Promise<T>;
  }
}
