export function Router(): JSX.Element;

export type ActionData = {
  actionId: string;
  data: any;
};

export type RouterState = {
  key: string;
  href: string;
  isInitial: boolean;
  actionData?: ActionData;
};

export function createRouterState(href: string): RouterState;
export function changeRouterState(href: string, key: string): RouterState;
export function changeRouterStateForAction(
  href: string,
  key: string,
  actionData: ActionData
): RouterState;

export function addNavigation(
  setRouter: (
    router: RouterState | ((router: RouterState) => RouterState)
  ) => void
): void;
jsxs;
export function navigate(href: string): void;
export function submitForm(actionData: ActionData): void;

declare global {
  interface Window {
    __rsc: any;
    __rscNav: (pathname: string) => void;
    __rscAction: (formAction: string, formData: FormData | undefined) => void;
  }
}
