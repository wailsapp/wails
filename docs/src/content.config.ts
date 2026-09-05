import { defineCollection } from "astro:content";
import { glob } from "astro/loaders";
import { docsSchema } from "@astrojs/starlight/schema";
import { blogSchema } from "starlight-blog/schema";

export const collections = {
  docs: defineCollection({
    // Translation sources stay in the repository until their pre-GA refresh,
    // but English is the only documentation published during the beta.
    loader: glob({
      base: "./src/content/docs",
      pattern: [
        "**/[^_]*.{markdown,mdown,mkdn,mkd,mdwn,md,mdx}",
        "!{de,fr,id,ja,ko,pt,ru,zh-cn,zh-tw}/**",
      ],
    }),
    schema: docsSchema({ extend: (context) => blogSchema(context) }),
  }),
};
