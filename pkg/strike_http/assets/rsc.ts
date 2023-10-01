import { ActionData } from "./router";

export type RemotePromise = {
  id: string;
  promise: Promise<any>;
  resolve: (value: any) => void;
  reject: (reason: any) => void;
};

export declare function createRemotePromise(id: string): RemotePromise;

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

export declare function RscComponent({
  isInitial,
  url,
  routerKey,
  actionData,
}: {
  isInitial: boolean;
  url: string;
  routerKey: string;
  actionData: ActionData;
}): React.ReactNode;
