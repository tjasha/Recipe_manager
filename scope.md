# Recipe manager application

Recipe manager applications web application that saves and categorizes the recipes. It also helps users prepare the shopping list. 
The goal is to develop a web application for managing recipes. This application will allow users to save, categorize, and organize their recipes, as well as facilitate the creation of shopping lists.
This project serves as a practical example for learning the Go programming language.

## Personas 

### Visitor
A user in this role has read-only access to the application, allowing them to view and search content but prohibiting the creation of new items. They are, however, able to register and transition to the "Chef" role.

### Chef
A registered user, the Chef, is able to perform various actions related to recipes, including creation, editing, categorization, and grading. Additionally, the Chef can favorite recipes and manage shopping lists by creating and editing them.

### System Admin 
The Admin role is granted complete application access. This includes the authority to manage all recipes across all users and the ability to remove users when required.

## Features: 
- User Authentication & Account Management
- Log in and Log out
- Register
- Change Password
- Recipe Management & Viewing
- Read Recipes
- Filter Recipes (by ingredients, category, cooking time)
- Adjust Recipe Portion Size
- Create, Edit, and Delete Recipes
- Add and Remove Recipes from Favourites
- Add and Remove Recipes from Collections
- Show and Adjust Cooking Steps
- List & Communication Features
- Create, Edit, and Delete Shopping Lists
- Export Recipe
- Receive Recipe via Email
- Exporting Shopping list
- Receive Shopping List via Email

## Workflows

As a **Visitor**, I should be able to:

- View all publicly available recipes.
- Search and filter recipes (by ingredients, category, cooking time).
- View a recipe's full details, including ingredients and cooking steps.
- Adjust the displayed recipe portion size when viewing a recipe.
- Register to transition to the "Chef" role.

As a **Chef**, I should be able to:

- Perform all actions available to a Visitor.
- Log in and Log out.
- View and manage my account details 
- Change my password.
- Create new recipes, including ingredients, steps, categories, and cooking time.
- Edit and Delete my own recipes.
- Add or remove recipes from my Favourites list.
- Organize recipes into Collections.
- Create, Edit, and Delete Shopping Lists.
- Export Recipe 
- Send a Shopping List to myself via Email.
- Export Shopping List.
- Send a specific Recipe to myself via Email.


As an **Admin**, I should be able to:

- Perform all actions available to a Chef.
- Edit and Delete any recipe created by any user.
- Manage users, including the ability to remove users from the application.
