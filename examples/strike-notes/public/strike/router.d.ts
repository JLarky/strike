export function Router(): JSX.Element;

export type RouterState = {
  key: string;
  href: string;
  isInitial: boolean;
};

export function createRouterState(href: string): RouterState;
export function changeRouterState(href: string, key: string): RouterState;

export function addNavigation(setRouter: (router: RouterState) => void): void;
jsxs;
export function navigate(href: string): void;

declare global {
  interface Window {
    __rsc: any;
    __rscNav: (pathname: string) => void;
  }
}
