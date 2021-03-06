import { User } from './accounting/user/user.class';
import { ProductService } from './product/product.service';
import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Router, NavigationStart } from '@angular/router';
import { AppService } from './app.service';
import { AuthService } from './accounting/auth.service';
import { Product } from './product/product';
import {
  trigger,
  state,
  style,
  animate
} from '@angular/animations';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/distinct';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass']
})
export class AppComponent implements OnInit {
  selectedProduct: Observable<Product>;
  authState: Observable<User>;
  addProductActive: Observable<boolean>;

  constructor(
    private router: Router,
    private $appService: AppService,
    private $authService: AuthService,
    private $productService: ProductService
  ) {}

  ngOnInit() {
    this.authState = this.$authService.fetchSession();
    this.addProductActive = this.$productService.fetchAddProductActive();
    this.addSubscriptions();
  }

  addSubscriptions() {
    this.$appService.watchForLoginConponent();
    this.selectedProduct = this.$productService.fetchSelectedProduct();
  }

  logout() {
    this.$authService.logout();
  }

  removeProduct(p: Product) {
    this.$productService.removeProduct(p);
  }
}
