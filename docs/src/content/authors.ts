import type { StarlightBlogUserConfig } from "starlight-blog";

type Authors = NonNullable<StarlightBlogUserConfig>["authors"];
export const authors: Authors = {
  leaanthony: {
    name: "Lea Anthony",
    title: "Maintainer of Wails",
    url: "https://github.com/leaanthony",
    picture: "https://github.com/leaanthony.png",
  },
  misitebao: {
    name: "Misite Bao",
    title: "Architect",
    url: "https://github.com/misitebao",
    picture: "https://github.com/misitebao.png",
  },
};
