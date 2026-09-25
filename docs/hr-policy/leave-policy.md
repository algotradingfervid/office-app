# Leave Policy (DRAFT — sample, to be tailored)

> **Status:** Dummy policy based on common practice in Indian IT companies (researched Sep 2026;
> aligned with the app design `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` on
> 2026-09-25). It is a starting point, not legal advice. Leave laws differ by state (Shops &
> Establishments Act of the state where the company is registered) and by the Labour Codes in force
> from 21 Nov 2025 (maternity benefit now sits in Chapter VI of the Code on Social Security, 2020) —
> have HR/legal confirm before adoption. Every number marked `[setting]` will be an admin-editable
> value in the app.

---

## 1. Purpose and scope

This policy explains the types of leave available to employees, how leave is earned, applied for,
approved, carried forward and encashed. It applies to all permanent employees. Interns and contract
staff follow their contract terms unless stated otherwise.

## 2. Leave year

- The leave year runs **1 April to 31 March**.
- Balances are carried forward or lapse on the first day of the leave year (see §4.3 and §5).
- A single leave application cannot cross the leave-year boundary. Leave from, say, 30 March to
  2 April is applied as two applications (30–31 March and 1–2 April), each charged to its own year.

## 3. Leave types at a glance

| Code | Leave type | Entitlement | How it is credited / counted | Carry forward | Encashable |
|------|-----------|---------------|--------------------|---------------|------------|
| CL | Casual Leave | 12 days per year `[setting]` | Monthly: 1 day on the 1st of each month `[setting]` | No — lapses on 31 Mar | No |
| SL | Sick Leave | 8 days per year `[setting]` | Full quota at start of year (pro-rata for joiners) | No — lapses on 31 Mar | No |
| EL | Earned / Privilege Leave | 18 days per year `[setting]` | Monthly: 1.5 days on the 1st of each month | Yes, max 30 days total balance at year end `[setting]` | Yes — at exit, and optionally above cap |
| ML | Maternity Leave | 26 weeks (182 **calendar** days) | On request, per event | — | No |
| PL | Paternity Leave | 5 working days `[setting]` | On request, per event | — | No |
| BL | Bereavement Leave | 3 working days `[setting]` | On request, per event | — | No |
| MRL | Marriage Leave | 3 working days `[setting]` | Once during employment | — | No |
| CO | Compensatory Off | As earned | **Not yet available** — will be introduced with the Comp-off claim form (§4.8) | — | No |
| LOP | Loss of Pay | Max 15 days per year `[setting]` | Used only when no paid balance is available | — | — |

## 4. Rules for each leave type

### 4.1 Casual Leave (CL)
- For short, unplanned personal needs.
- Credited **1 day per month** `[setting]`, in advance on the 1st, so an employee can only use what has
  been credited so far (e.g., by the end of June, at most 3 days).
- Maximum **2 consecutive working days** `[setting]`. Casual leave applications that are next to each
  other (with only weekly offs or holidays between them) are counted together. Longer absences must
  use EL.
- Apply at least **1 day in advance** `[setting]` where possible. Same-day applications are accepted in
  emergencies; the approver sees that the notice was short.
- Cannot be applied for a past date `[setting: 0 days backdating]`.

### 4.2 Sick Leave (SL)
- For illness or medical appointments. May be applied on the day or after returning to work, within
  **7 days** of the first day of absence `[setting]`.
- A **medical certificate** must be attached to the application for more than **2 consecutive working
  days** `[setting]`.
- Medical certificates are confidential: they are visible only to the employee, the approvers of that
  application and HR.

### 4.3 Earned / Privilege Leave (EL)
- For planned leave such as vacations.
- Apply at least **7 days in advance** `[setting]`. Shorter notice is not rejected automatically, but the
  approver sees that the notice was short and may reject.
- Cannot be applied for a past date `[setting: 0 days backdating]`.
- Unused EL is carried forward, but the balance carried into the new year cannot exceed **30 days**
  `[setting]`; days above the cap lapse (or are encashed, if the company chooses `[setting]`).
  The cap is applied only at year end — during the year the balance may go above it.
- Unused EL balance is **paid out at exit** (resignation, retirement, termination).

### 4.4 Maternity Leave (ML)
- **26 weeks paid** (182 calendar days, including Sundays and holidays) for a woman with fewer than two
  surviving children; **12 weeks** for a woman with two or more surviving children (Code on Social
  Security, 2020, Chapter VI — formerly the Maternity Benefit Act, as amended 2017).
- Eligibility: at least **80 days** worked in the 12 months before the expected delivery date.
- Up to 8 weeks may be taken before the expected delivery date.
- Eligibility and the applicable number of weeks are confirmed by the approver/HR; the app counts the
  days and allows up to 182.

### 4.5 Paternity Leave (PL)
- Not mandated by law for the private sector; this is a company benefit.
- **5 working days**, to be taken within **3 months** of the child's birth or adoption (confirmed by
  the approver).

### 4.6 Bereavement Leave (BL)
- **3 working days** on the death of an immediate family member (spouse, child, parent, sibling,
  parent-in-law).

### 4.7 Marriage Leave (MRL)
- **3 working days**, once during employment, for the employee's own marriage.

### 4.8 Compensatory Off (CO)
- Not available in the first release of the app. It will be introduced with a separate **Comp-off
  claim** form, under these intended rules:
  - Earned when an employee works a full day on a weekly off or public holiday **with prior
    manager approval**.
  - Must be used within **60 days** of being earned `[setting]`, otherwise it lapses.
- Until then, work on a weekly off or holiday is handled by HR as agreed with the manager.

### 4.9 Loss of Pay (LOP)
- Applies when the employee has no balance in the requested leave type, or for unauthorised absence.
- Salary is deducted for LOP days. HR exports LOP days per employee per month for payroll.
- LOP should be used only when no paid leave is available; if an employee applies for LOP while CL or
  EL balance is available, the approver is told.
- LOP is **capped at 15 days per leave year** `[setting]` (approved and pending LOP together). The app
  will not accept an LOP application beyond the cap; the employee must contact HR, and only the HR admin
  can grant an exception by recording the leave with a reason.
- A single LOP application may not exceed **5 consecutive working days** `[setting]`.
- Unauthorised absence is recorded as LOP by HR.

## 5. How leave is counted

- **Working week:** Monday to Saturday; Sundays and the **2nd and 4th Saturday** of each month are
  weekly offs `[setting]`. The same weekly offs and holiday list apply to all employees.
- **Half days** are allowed for CL, SL and EL `[setting]`:
  - Single day: full day, first half or second half.
  - Several days: the first day may start at the second half and the last day may end at the first
    half.
  - A half day counts as 0.5 day.
- **Weekends and public holidays** falling inside a leave period are **not** counted as leave
  `[setting]`, except for Maternity Leave, which is counted in calendar days.
- **Sandwich rule:** not applied — a weekend between two leave days is not counted `[setting]`.
- A leave application must **start and end on a working day** (except Maternity Leave).
- **No overlapping applications:** an employee cannot apply for dates already covered by another
  pending or approved application (a first-half and a second-half on the same day are allowed).
- **Negative balance:** not allowed. The app blocks an application that exceeds the available balance;
  the employee must apply the excess days separately as LOP (subject to the LOP cap).
- **Available balance** = balance minus days already applied for and still pending, so the same days
  cannot be booked twice.
- **Holidays added later:** the days are counted again at final approval. If a holiday has been added
  since the application, the corrected count is deducted.

## 5A. Settings available for every leave type

So the policy can change without code changes, each leave type has these admin-editable settings:

| Setting | Example (CL) |
|---------|--------------|
| Days per year (or "no quota") | 12 |
| Credit method: full quota at year start / monthly / on request per event | Monthly |
| Count working days or calendar days | Working days |
| Maximum consecutive days per application | 2 |
| Maximum days per year (hard cap, even if no quota) | 12 |
| Maximum times during employment | — |
| Minimum advance notice (days) — shown to the approver as a warning, not a block | 1 |
| Maximum days an application may be backdated | 0 |
| Half day allowed | Yes |
| Attachment required after N consecutive days | — |
| Carry forward cap (0 = lapses) | 0 |
| Encashable at exit | No |
| Available during probation | Yes |
| Effective from date | 1 Apr 2027 |

## 6. New joiners, probation and exit

- **Pro-rata credit:** SL is credited in proportion to the months remaining in the leave year,
  including the month of joining, rounded down to the nearest half day (e.g., joining in October →
  6/12 of the annual quota). CL and EL are credited in advance on the 1st of each month to employees who
  have joined on or before that date, so a mid-month joiner's first credit is on the next 1st.
- **During probation** (first 6 months `[setting]`): CL and SL may be used; EL accrues but can be
  used only after the probation end date recorded by HR `[setting]`.
- **At exit:** HR deactivates the employee's account; pending applications are cancelled. Unused EL is
  encashed; CL and SL lapse. Leave taken in excess of the pro-rata entitlement is recovered from final
  settlement.
- **Records of ex-employees:** leave records are kept; medical certificates are deleted after a
  retention period set by HR `[setting]`.

## 7. Public holidays

- The company publishes a holiday list of **10 days** `[setting]` before the start of each leave year,
  including the national holidays (Republic Day, Independence Day, Gandhi Jayanti).
- The holiday list is maintained in the app and is used to calculate leave days.

## 8. Applying for and approving leave

1. The employee submits a leave application in the app: leave type, from/to dates, half-day option,
   reason, and attachment (if required, e.g., medical certificate).
2. The employee selects the **first approver** from the list of users with approval rights (not
   themselves).
3. Each approver can **Approve & forward** to another approver, **Approve as final**, or **Reject**
   with a comment. An application can be forwarded at most 5 times.
4. Leave balance is **deducted only on final approval**. While pending, the days are shown as
   "on hold" so they cannot be double-booked.
5. An application cannot be edited after submission; cancel it and apply again.
6. **Cancellation:**
   - A pending application can be cancelled by the employee at any time.
   - Approved leave that has not yet started: the employee asks to cancel, and the final approver
     confirms or declines. On confirmation the balance is restored.
   - Leave that has started or is in the past (including returning early): only HR can cancel it. HR
     then records the days actually taken.
7. If an approver is absent or leaves, HR reassigns their pending applications to another approver.
8. HR may **record leave on an employee's behalf** (e.g., unauthorised absence as LOP, an LOP
   exception beyond the cap, or an employee without access to the app). A reason is required and the
   record shows that HR entered it.
9. Every action is recorded in the application's history and cannot be edited.

## 9. Policy administration

- HR owns this policy and may revise it. Changes apply from the stated effective date to applications
  submitted on or after that date; they do not alter balances, or applications already submitted or
  approved.
- Values marked `[setting]` are configured by the HR admin in the application; every change is
  recorded in an audit log.

---

### Sources consulted
- Payoneer — Leave policy in India 2026: https://www.payoneer.com/en-in/resources/workforce-management/leave-policy/india/
- AYP Group — India Leave Policy Employer Guide 2026: https://ayp-group.com/blog/leave-policy-in-india
- WageIndicator — Earned & Casual Leave in India: https://wageindicator.org/en-in/work-in-india/labour-law/annual-leave-and-holidays/earned-casual-leave-in-india/
- SalaryBox — EL/CL/SL rules, accumulation & encashment 2026: https://salarybox.in/blog/earned-leave-casual-leave-sick-leave-in-india-rules-accumulation-encashment-explained-2026/
- greytHR — Maternity Leave in India 2026: https://www.greythr.com/blog/maternity-leave/
- Wisemonk — Paternity leave in India 2026: https://www.wisemonk.io/blogs/paternity-leave-in-india
- Maternity benefit includes Sundays and paid holidays (B. Shah v. Presiding Officer, 1977): https://karmamgmt.com/blog/maternity-benefit-includes-sundays-paid-holidays-india
- greytHR — leave spanning the year end: https://www.greythr.com/employee-portal/answers/40768982/
