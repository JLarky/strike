export type RemotePromise = {
  id: string;
  promise: Promise<any>;
  resolve: (value: any) => void;
  reject: (reason: any) => void;
};

export function createRemotePromise(id: string): RemotePromise;

export function remotePromiseFromCtx(ctx: CTX, id: string): RemotePromise;

type CTX = { promises: Map<string, RemotePromise> };

export function parseModelString(
  ctx: CTX,
  parent: { [key: string]: string | null | Symbol },
  key: string,
  value: string
): any;

function promisify(obj: { [key: string]: any }, promise): void {}
