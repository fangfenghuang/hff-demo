package com.hfftest.java8demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Java8demoApplication {

    public static void main(String[] args) {
        SpringApplication.run(Java8demoApplication.class, args);
        CyclePrint.cyclePrint();
    }

}
