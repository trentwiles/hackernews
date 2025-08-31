import type { hr } from "date-fns/locale";
import { Underline } from "lucide-react";
import React from "react";
import ReactMarkdown from "react-markdown";
import rehypeSanitize from "rehype-sanitize";

interface props {
  content: string;
}

export default function SafeMarkdown(props: props) {
  const customComps = {
    h1(props: any) {
      const { children, ...rest } = props;
      return (
        <h1
          className="scroll-m-20 text-3xl font-extrabold tracking-tight text-balance"
          {...rest}
        >
          {children}
        </h1>
      );
    },

    h2(props: any) {
      const { children, ...rest } = props;
      return (
        <h1
          className="scroll-m-20 text-2xl font-extrabold tracking-tight text-balance"
          {...rest}
        >
          {children}
        </h1>
      );
    },

    h3(props: any) {
      const { children, ...rest } = props;
      return (
        <h1
          className="scroll-m-20 text-1xl font-extrabold tracking-tight text-balance"
          {...rest}
        >
          {children}
        </h1>
      );
    },

    h4(props: any) {
      const { children, ...rest } = props;
      return (
        <h1
          className="scroll-m-20 text-xl font-extrabold tracking-tight text-balance"
          {...rest}
        >
          {children}
        </h1>
      );
    },

    a(props: any) {
      const { children, href, ...rest } = props;
      return (
        <a
          style={{textDecoration: 'underline', color: 'blue'}}
          target="_blank"
          href={"/urlCheck?q=" + href}
          {...rest}
        >
          {children}
        </a>
      );
    },
  };

  return (
    <ReactMarkdown rehypePlugins={[rehypeSanitize]} components={customComps}>
      {props.content}
    </ReactMarkdown>
  );
}
