import { ActionData } from "./router";

export type RemotePromise = {
  id: string;
  promise: Promise<any>;
  resolve: (value: any) => void;
  reject: (reason: any) => void;
};

export declare function chunkToJSX(ctx: CTX, str: string): any;

export declare function createRemotePromise(id: string): RemotePromise;

export declare function fetchChunksPromise(
  href: string
): Promise<React.ReactNode>;

export declare function fetchFromActionPromise(
  href: string,
  actionData: ActionData
): Promise<React.ReactNode>;

export declare function remotePromiseFromCtx(
  ctx: CTX,
  id: string
): RemotePromise;

type CTX = { promises: Map<string, RemotePromise> };

export declare function parseModelString(
  ctx: CTX,
  parent: { [key: string]: string | null | Symbol },
  key: string,
  value: string
): any;

export declare function promisify<T>(
  obj: { [key: string]: any },
  promise: Promise<T>
): void;
export declare function actionify(
  obj: { [key: string]: any },
  actionId: string
): void;

export type RscComponentProps = {
  isInitial: boolean;
  url: string;
  urlPromise: Promise<unknown> | undefined;
  routerKey: string;
  actionData: ActionData | undefined;
  actionPromise: Promise<unknown> | undefined;
};

export declare function RscComponent({
  isInitial,
  url,
  urlPromise,
  routerKey,
  actionData,
  actionPromise,
}: RscComponentProps): React.ReactNode;
