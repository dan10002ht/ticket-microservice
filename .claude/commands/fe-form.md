# Create Form Component

You are a principal frontend developer. Generate a form with react-hook-form + zod validation.

## Instructions

1. Check the backend validation rules in `gateway/src/middlewares/validationMiddleware.js` to understand what fields are required and their constraints
2. Create the zod schema in `frontend/src/lib/validators/`
3. Create the form component with react-hook-form

## Form Pattern

```tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

// 1. Define schema
const formSchema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type FormValues = z.infer<typeof formSchema>;

// 2. Define props
interface LoginFormProps {
  onSubmit: (values: FormValues) => void;
  isLoading?: boolean;
}

// 3. Component
export function LoginForm({ onSubmit, isLoading }: LoginFormProps) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="you@example.com" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Loading..." : "Submit"}
        </Button>
      </form>
    </Form>
  );
}
```

## Validator Pattern

```typescript
// lib/validators/auth.ts
import { z } from "zod";

export const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
});

export const registerSchema = z.object({
  firstName: z.string().min(1, "First name is required"),
  lastName: z.string().min(1, "Last name is required"),
  email: z.string().email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
```

## Backend Validation Reference

Check `gateway/src/middlewares/validationMiddleware.js` for field rules:
- `validateRegistration` - email, password, firstName, lastName
- `validateLogin` - email, password
- `validateEvent` - name, description, start_date, end_date, venue_name, venue_address, venue_city, venue_country, venue_capacity
- `validateBooking` - eventId, ticketQuantity
- `validatePayment` - bookingId, paymentMethod, amount
- `validateTicketCreate` - eventId, ticketTypeId, price
- `validateUserProfileCreate` - firstName, lastName, email, phone
- `validateUserAddress` - street, city, state, postalCode, country

**Note:** Gateway auto-converts camelCase (frontend) ↔ snake_case (backend). Use snake_case in zod schemas to match the API.

## shadcn/ui Form Components

Use these shadcn/ui components for forms:
- `Form`, `FormField`, `FormItem`, `FormLabel`, `FormControl`, `FormMessage` - Form wrapper
- `Input` - Text input
- `Textarea` - Multi-line text
- `Select`, `SelectTrigger`, `SelectContent`, `SelectItem` - Dropdown
- `Checkbox` - Checkbox
- `RadioGroup`, `RadioGroupItem` - Radio buttons
- `DatePicker` - Date selection
- `Switch` - Toggle

## Rules

- Zod schemas in `lib/validators/{domain}.ts`, exported as both schema and inferred type
- Form components receive `onSubmit` callback and `isLoading` prop
- Form components do NOT call API directly - parent page handles that via mutation hooks
- Always provide default values in useForm
- Use `<Label>` + `<Input>` + inline error `<p>` for simple forms (actual pattern used in this project)
- Use shadcn/ui `Form`/`FormField` wrapper only for complex forms
- Use toast (sonner) for form-level/API errors
- Disable submit button when `isLoading`

## User Input

$ARGUMENTS
