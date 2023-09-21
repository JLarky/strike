/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

"use client";

import { useTransition } from "react";

export default function EditButton({ noteId, children }) {
  const [isPending, startTransition] = useTransition();
  const isDraft = noteId == null;
  return (
    <button
      className={[
        "edit-button",
        isDraft ? "edit-button--solid" : "edit-button--outline",
      ].join(" ")}
      disabled={isPending}
      onClick={() => {
        startTransition(() => {
          __rscNav(isDraft ? "/edit" : `/edit/${noteId}`);
        });
      }}
      role="menuitem"
    >
      {children}
    </button>
  );
}
