import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

const routes: Routes = [];

@NgModule({
  imports: [
    RouterModule.forRoot(routes,{useHash:true})
  ],
  exports: [RouterModule]
})

export class AppRoutingModule { }
