export function Router(): JSX.Element;

export type RouterState = {
  key: string;
  path: string;
  isInitial: boolean;
};

export function createRouterState(path: string): RouterState;
export function changeRouterState(path: string, key: string): RouterState;

export function addNavigation(setRouter: (router: RouterState) => void): void;
jsxs;
export function navigate(path: string): void;

declare global {
  interface Window {
    __rsc: any;
    __rscNav: (pathname: string) => void;
  }
}
