# Wails Enhancement Proposal (WEP) 

## Title

**Author**: Lea Anthony
**Created**: 2024-05-19

## Summary

The move to services instead of bound structs is a great move forward for the project. This allows us to have a more flexible and powerful tooling system. This proposal outlines a standard way of creating, managing, and using services in Wails applications.

## Motivation

The current system of using bound structs has several limitations, including difficulties in creating portable services, lack of standardized initialisation and de-initialisation patterns, and the absence of CLI tooling to manage services efficiently. By moving to a service-based architecture, we can address these issues and provide a more robust and developer-friendly environment.

## Detailed Design

#### Technical Details

- **Service Architecture**: Services will be modular and self-contained, allowing for easy portability and reuse across different projects.
- **Standardized Patterns**: Implement standardized initialisation and de-initialisation patterns to ensure consistency and reliability.
- **CLI Tooling**: Introduce `wails3 service <command>` to provide Laravel-like tooling for managing services. This will include commands for creating and deleting services, as well as creating and deleting service methods.

#### Implementation Steps

1. **Define Service Interface**: A standard interface is created that all services must implement:
    ```go
    type Service interface {
         Init() error
         Shutdown() error
    }
    ```
2. **Define Service File Layout**: A standard file layout is defined for services:
```
<project>/
    └── services/
        └── myservice/
            ├── myservice.go
            ├── service.yaml
            ├── go.mod
            ├── README.md
            ├── method1.go
            └── method2.go
```

`myservice.go` contains the service implementation. This will be named according to the service name.

`service.yaml` contains metadata about the service.

`go.mod` is the standard Go module file.

`README.md` provides documentation for the service.

`method1.go` & `method2.go` represent the service methods.


3. **Develop CLI Tooling**: Implement the `wails3 service <command>` CLI tool with commands for managing services. Initially, we will only implement the `create` command: `wails3 service create <service-name>`. This will create a new service which will involve the following steps:
      - Create the service directory structure.
      - Generate the `myservice.go` file with the service interface implementation.
      - Generate the `service.yaml` file with metadata about the service.
      - Generate the `README.md` file with documentation for the service.
      - Generate the `go.mod` file for the service.
      
4. **Documentation**: Update documentation to reflect the new service-based architecture and CLI tooling.

#### Potential Impact on Existing Functionality

As this is a new feature, there should be minimal impact on existing functionality. However, existing projects will need to be refactored to adopt the new service architecture.

## Pros/Cons

#### Pros

- **Modularity**: Services are self-contained and portable, making it easier to reuse and share code.
- **Consistency**: Standardized initialisation and de-initialisation patterns ensure consistent behavior across services.
- **Developer Productivity**: CLI tooling simplifies service management, reducing the time and effort required for common tasks.

#### Cons

- **Migration Effort**: Existing projects will need to invest time in refactoring their code to adopt the new architecture.
- **Learning Curve**: Developers will need to learn the new service patterns and CLI commands. I don't believe this is a significant issue as the benefits outweigh the learning curve.

## Alternatives Considered

- Not providing tooling. It's perfectly acceptable for developers to manage services manually, but tooling can greatly improve productivity.

## Backwards Compatibility

As this will be a new feature, there should be minimal impact on existing projects. 

## Test Plan

- **Unit Tests**: The tooling will utilise unit tests to determine if the service has been created correctly.

## Reference Implementation

There is currently no reference implementation, though the functionality would be very similar to the current `wails3 init` command.

## Maintenance Plan

- This feature will be maintained by the Wails maintainers.

## Conclusion

The move to a service-based architecture offers significant benefits, and providing tooling to manage services will further enhance developer productivity. By adopting this proposal, we can create a  robust foundation for managing services in Wails applications.

