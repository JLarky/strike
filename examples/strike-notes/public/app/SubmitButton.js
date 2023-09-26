import { createElement, useEffect, useRef } from "react";
import { experimental_useFormStatus as useFormStatus } from "react-dom";
import { jsx } from "react/jsx-runtime";

export function SubmitButton(props) {
  const { pending } = useFormStatus();
  console.log(pending, props.myAct);
  const child = props.children[0];
  const ref = useRef();
  useEffect(() => {
    const form = ref.current?.form;
    if (form) {
      form.addEventListener("submit", (e) => {
        e.preventDefault();
        props.myAct.action(new FormData(form));
      });
    }
  }, [ref.current]);
  return createElement(child.type, { ...child.props, ref });
}
