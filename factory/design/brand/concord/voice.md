# Voice — Concord

Placet speaks like a good HR officer: polite, plain and exact. It says what happened, to whom, and the
numbers; it never jokes about someone's leave.

**Adjectives:** polite, plain, exact.

| We say | We don't say |
|---|---|
| "Approved. 2 days of casual leave deducted." | "Woohoo! Your leave is approved 🎉" |
| "Casual leave allows at most 2 consecutive days." | "Error: validation failed for field leave_days." |
| "Sent to Ravi Kumar for approval." | "Your request has been forwarded to the next approver in the workflow." |

Rules: exact numbers, in the mono face (`2 days`, `REQ-0142`). No exclamation marks. Name the person when we
know them. Say what happens next. Leave types in full on first mention ("casual leave"), short forms (CL, EL)
in tables where staff already use them.

## Five real strings
- **Welcome:** "Good morning, Asha. You have 4 days of casual leave available."
- **Empty state:** "No requests are waiting for your approval."
- **Error:** "Casual leave allows at most 2 consecutive days. Shorten the dates or choose earned leave."
- **Success toast:** "Approved. 2 days of casual leave deducted from Asha's balance."
- **Destructive confirm:** "Cancel this leave? The 2 days return to your casual leave balance, and the cancellation is recorded in the history." Buttons: "Cancel leave" / "Keep leave".
