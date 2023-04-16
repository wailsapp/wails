// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

export namespace main {
  
  export type PersonSource = Partial<{
    name: string;
    address: services.Address;
  }>
  
  export class Person {
    name: string;
    address: services.Address;
    
    static createFrom(source: string | PersonSource = {}) {
      return new Person(source);
    }

    constructor(source: string | PersonSource = {}) {
      if ('string' === typeof source) {
        source = JSON.parse(source);
      }

      this.name = source['name'];
      this.address = services.Address.createFrom(source['address']);
      
    }
  }
  
}

export namespace services {
  
  export type AddressSource = Partial<{
    street: string;
    state: string;
    country: string;
  }>
  
  export class Address {
    street: string;
    state: string;
    country: string;
    
    static createFrom(source: string | AddressSource = {}) {
      return new Address(source);
    }

    constructor(source: string | AddressSource = {}) {
      if ('string' === typeof source) {
        source = JSON.parse(source);
      }

      this.street = source['street'];
      this.state = source['state'];
      this.country = source['country'];
      
    }
  }
  
}
