// app/EditButton.jsx
import {useTransition} from "react";
import {
jsx
} from "react/jsx-runtime";
function EditButton({ noteId, children }) {
  const [isPending, startTransition] = useTransition();
  const isDraft = noteId == null;
  return jsx("button", {
    className: [
      "edit-button",
      isDraft ? "edit-button--solid" : "edit-button--outline"
    ].join(" "),
    disabled: isPending,
    onClick: () => {
      startTransition(() => {
        __rscNav(isDraft ? "/edit" : `/edit/${noteId}`);
      });
    },
    role: "menuitem",
    children
  });
}
export {
  EditButton as default
};
