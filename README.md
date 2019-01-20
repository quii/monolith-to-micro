# Monolith to Micro

Starting microservices I feel is a waste of time and a violation of YAGNI. That's not to say just build a big monolith! You should start with a monolith to start building your product, iterating and evolving it listening to feedback. Eventually you'll start to "see" the services you'll want to split out, based on real experience of the system rather than architect trying to guess it. 

## General ideas

- To keep running things consistent use docker-compose, even for the first iteration. That's not too much overhead and will make things gentler as we add new things.
- Use gRPC to split things out.
- Keep it as a command line app just to minimise things.

## The problem

We want to know what to make for dinner!

There will be some kind of idea of what ingredients are in the house and what their expiration dates are. We'll also have a recipe book to derive meals from, which eventually we should be able to filter by 

### How to break the problem down

1. **Hello, world**. Command line app running through docker-compose that prints hello, world
2. **Hard-coded ingredients to use**. Print out a list of ingredients that are available ordered by the expiration date
3. **Manage ingredients (add, delete)**
4. **Find meals** from a hardcoded list of recipes and print them instead, based on available ingredients
5. **Return meals that dont have all ingredients** and list them
6. **Manage ingredients**

At this point, we'll think about splitting into different gRPC services 