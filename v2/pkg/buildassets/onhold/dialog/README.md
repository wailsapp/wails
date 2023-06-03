## Dialog

NOTE: Currently, this is a Mac only feature.

Place any PNG file in this directory to be able to use them in message dialogs.
The files should have names in the following format:
`name[-(light|dark)][2x].png`

Examples:

- `mypic.png` - Standard definition icon with ID `mypic`
- `mypic-light.png` - Standard definition icon with ID `mypic`, used when system
  theme is light
- `mypic-dark.png` - Standard definition icon with ID `mypic`, used when system
  theme is dark
- `mypic2x.png` - High definition icon with ID `mypic`
- `mypic-light2x.png` - High definition icon with ID `mypic`, used when system
  theme is light
- `mypic-dark2x.png` - High definition icon with ID `mypic`, used when system
  theme is dark

### Order of preference

Icons are selected with the following order of preference:

For High Definition displays:

- name-(theme)2x.png
- name2x.png
- name-(theme).png
- name.png

For Standard Definition displays:

- name-(theme).png
- name.png
