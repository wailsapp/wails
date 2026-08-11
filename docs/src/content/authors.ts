import type { StarlightBlogUserConfig } from "starlight-blog";

// starlight-blog v0.28 accepts either one blog config or an array of them.
// Authors live on a single config, so narrow the union before indexing it.
type BlogConfig = Exclude<NonNullable<StarlightBlogUserConfig>, readonly unknown[]>;
type Authors = BlogConfig["authors"];
export const authors: Authors = {
  leaanthony: {
    name: "Lea Anthony",
    title: "Maintainer of Wails",
    url: "https://github.com/leaanthony",
    picture: "https://github.com/leaanthony.png",
  },
};
