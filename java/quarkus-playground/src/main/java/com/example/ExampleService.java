package com.example;

import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class ExampleService {
    public String greeting(String name) {
        return "hello " + name;
    }
}
