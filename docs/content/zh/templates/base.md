---
title: Base 模板 and Blocks
linktitle:
description: The base and block constructs allow you to define the outer shell of your master templates (i.e., the chrome of the page).
godocref: https://golang.org/pkg/text/template/#example_Template_block
date: 2017-02-01
publishdate: 2018-08-11
lastmod: 2017-02-01
categories: [templates,fundamentals]
keywords: [blocks,base]
menu:
  docs:
    parent: "templates"
    weight: 20
weight: 20
sections_weight: 20
draft: false
aliases: [/templates/blocks/,/templates/base-templates-and-blocks/]
toc: true
---

The `block` keyword allows you to define the outer shell of your pages' one or more master template(s) and then fill in or override portions as necessary.

{{< youtube QVOMCYitLEc >}}

## Base Template Lookup Order

The [lookup order][lookup] for base templates is as follows:

1. `/layouts/section/<TYPE>-baseof.html`
2. `/themes/<THEME>/layouts/section/<TYPE>-baseof.html`
3. `/layouts/<TYPE>/baseof.html`
4. `/themes/<THEME>/layouts/<TYPE>/baseof.html`
5. `/layouts/section/baseof.html`
6. `/themes/<THEME>/layouts/section/baseof.html`
7. `/layouts/_default/<TYPE>-baseof.html`
8. `/themes/<THEME>/layouts/_default/<TYPE>-baseof.html`
9. `/layouts/_default/baseof.html`
10. `/themes/<THEME>/layouts/_default/baseof.html`

Variables are denoted by capitalized text set within `<>`. Note that Hugo's default behavior is for `type` to inherit from `section` unless otherwise specified.

### Example Base Template Lookup Order

As an example, let's assume your site is using a theme called "mytheme" when rendering the section list for a `posts` section. Hugo picks `layout/section/posts.html` as the template for [rendering the section][]. The `{{define}}` block in this template tells Hugo that the template is an extension of a base template.

Here is the lookup order for the `posts` base template:

1. `/layouts/section/posts-baseof.html`
2. `/themes/mytheme/layouts/section/posts-baseof.html`
3. `/layouts/posts/baseof.html`
4. `/themes/mytheme/layouts/posts/baseof.html`
5. `/layouts/section/baseof.html`
6. `/themes/mytheme/layouts/section/baseof.html`
7. `/layouts/_default/posts-baseof.html`
8. `/themes/mytheme/layouts/_default/posts-baseof.html`
9. `/layouts/_default/baseof.html`
10. `/themes/mytheme/layouts/_default/baseof.html`

## Define the Base Template

The following defines a simple base template at `_default/baseof.html`. As a default template, it is the shell from which all your pages will be rendered unless you specify another `*baseof.html` closer to the beginning of the lookup order.

{{< code file="layouts/_default/baseof.html" download="baseof.html" >}}
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{ block "title" . }}
      <!-- Blocks may include default content. -->
      {{ .Site.Title }}
    {{ end }}</title>
  </head>
  <body>
    <!-- Code that all your templates share, like a header -->
    {{ block "main" . }}
      <!-- The part of the page that begins to differ between templates -->
    {{ end }}
    {{ block "footer" . }}
    <!-- More shared code, perhaps a footer but that can be overridden if need be in  -->
    {{ end }}
  </body>
</html>
{{< /code >}}

## Override the Base Template

From the above base template, you can define a [default list template][hugolists]. The default list template will inherit all of the code defined above and can then implement its own `"main"` block from:

{{< code file="layouts/_default/list.html" download="list.html" >}}
{{ define "main" }}
  <h1>Posts</h1>
  {{ range .Pages }}
    <article>
      <h2>{{ .Title }}</h2>
      {{ .Content }}
    </article>
  {{ end }}
{{ end }}
{{< /code >}}

This replaces the contents of our (basically empty) "main" block with something useful for the list template. In this case, we didn't define a `"title"` block, so the contents from our base template remain unchanged in lists.

{{% warning %}}
Code that you put outside the block definitions *can* break your layout. This even includes HTML comments. For example:

```
<!-- Seemingly harmless HTML comment..that will break your layout at build -->
{{ define "main" }}
...your code here
{{ end }}
```
[See this thread from the Hugo discussion forums.](https://discourse.gohugo.io/t/baseof-html-block-templates-and-list-types-results-in-empty-pages/5612/6)
{{% /warning %}}

The following shows how you can override both the `"main"` and `"title"` block areas from the base template with code unique to your [default single page template][singletemplate]:

{{< code file="layouts/_default/single.html" download="single.html" >}}
{{ define "title" }}
  <!-- This will override the default value set in baseof.html; i.e., "{{.Site.Title}}" in the original example-->
  {{ .Title }} &ndash; {{ .Site.Title }}
{{ end }}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
{{< /code >}}

[hugolists]: /templates/lists
[lookup]: /templates/lookup-order/
[rendering the section]: /templates/section-templates/
[singletemplate]: /templates/single-page-templates/
