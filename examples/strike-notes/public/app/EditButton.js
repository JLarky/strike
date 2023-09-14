// app/EditButton.jsx
import {useTransition} from "react";
import {
jsx
} from "react/jsx-runtime";
function EditButton({ noteId, title, children }) {
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
    children: title
  });
}
export {
  EditButton as default
};
