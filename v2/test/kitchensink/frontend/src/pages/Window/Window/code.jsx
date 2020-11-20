import { Window } from '@wails/runtime';

function resize(imageWidth, imageHeight) {
  Window.SetSize(imageWidth, imageHeight);
  Window.Center();
}