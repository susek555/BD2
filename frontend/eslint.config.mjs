import { FlatCompat } from "@eslint/eslintrc";
import { dirname } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
    baseDirectory: __dirname,
});

const eslintConfig = [
  ...compat.config({
    extends: ["next/core-web-vitals", "next/typescript"],
    rules: {
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: '^_', // Ignore unused function arguments
          varsIgnorePattern: '^_', // Ignore unused variables
          caughtErrorsIgnorePattern: '^_', // Ignore unused catch clause parameters
          ignoreRestSiblings: true, //  Ignore unused properties when using object destructuring
        },
      ],
    }
  }),
];

export default eslintConfig;
