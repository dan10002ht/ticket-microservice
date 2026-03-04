import { z } from "zod";

export const profileSchema = z.object({
  firstName: z.string().min(1, "Vui lòng nhập họ"),
  lastName: z.string().min(1, "Vui lòng nhập tên"),
  phone: z
    .string()
    .regex(/^\+?[\d\s-]{9,15}$/, "Số điện thoại không hợp lệ")
    .optional()
    .or(z.literal("")),
});

export type ProfileInput = z.infer<typeof profileSchema>;
