# Blacklist App

Noona HQ example app that adds blacklisting functionality.

## How it works

### Setup phase

1. User installs app
2. App creates a customer group called **Blacklist**
3. App creates a webhook to track appointment (event) creation

### Daily workings

The app examines all appointments (events) that are created. If an event satisfies the following:

- Appointment is created through marketplace
- Attached customer is in **Blacklist** customer groups

the appointment is automatically declined.

### Use case

Users sometimes encounter individuals that could be labelled a "Bad Client". They frequently don't show up, they don't pay, spam appointment requests without intending to ever show up etc.

The user can now simply add this customer to the **Blacklist** customer group and they will be freed from future nuisance.
