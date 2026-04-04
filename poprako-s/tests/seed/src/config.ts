export const BASE_URL = "http://127.0.0.1:8080/api/v1";

export const SUPER_ADMIN_CREDS = {
  qq: "123456789",
  password: "123456",
};

// Unique suffix for this run — avoids unique-constraint collisions on repeated runs
const RUN_ID = Date.now() % 999_999;

// Role bitmasks
export const ROLE = {
  RAW_PROVIDER: 1,
  TRANSLATOR: 2,
  PROOFREADER: 4,
  TYPESETTER: 8,
  REVIEWER: 16,
  PUBLISHER: 32,
  ADMIN: 64,
} as const;

// Workflow transition strings
export const WORKFLOW = {
  UPLOAD_COMPLETE: "upload_complete",
  TRANSLATE_START: "translate_start",
  TRANSLATE_COMPLETE: "translate_complete",
  PROOFREAD_START: "proofread_start",
  PROOFREAD_COMPLETE: "proofread_complete",
  TYPESET_START: "typeset_start",
  TYPESET_COMPLETE: "typeset_complete",
  REVIEW_COMPLETE: "review_complete",
  PUBLISH_COMPLETE: "publish_complete",
} as const;

// Test user credentials — QQ is unique per run
export const TRANSLATOR_CREDS = {
  qq: `${20000 + RUN_ID}`,
  password: "test123",
  name: `Translator-${RUN_ID}`,
};

export const PROOFREADER_CREDS = {
  qq: `${30000 + RUN_ID}`,
  password: "test123",
  name: `Proofreader-${RUN_ID}`,
};
