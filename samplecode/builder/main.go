package main

import (
	"context"
	"fmt"
	"time"

	pdf "github.com/chinmay-sawant/gopdfsuit-client"
)

const (
	baseURL    = "http://localhost:8080"
	outputPath = "output_from_builder.pdf"
)

func main() {
	// Create a client
	client := pdf.NewClient(
		baseURL,
		pdf.WithTimeout(60*time.Second),
		pdf.WithMaxRetries(3),
	)

	ctx := context.Background()

	fmt.Println("=== Building Patient Registration Form and sending to endpoint ===")

	// 1. Build the document object
	doc := buildPatientForm()

	fmt.Printf("Successfully built document: %s\n", doc.Title.Text)
	fmt.Printf("Document contains %d tables\n", len(doc.Tables))

	// 2. Send the document
	fmt.Printf("Sending document to %s...\n", baseURL)
	err := client.SendAndSave(ctx, doc, outputPath)
	if err != nil {
		fmt.Printf("Note: Request failed (%v).\n", err)
		fmt.Println("Make sure the PDF generation service is running at", baseURL)
		return
	}

	fmt.Printf("Success! PDF saved to: %s\n", outputPath)
}

func buildPatientForm() *pdf.Document {
	// Common properties
	labelProps := "font1:9:100:left:1:1:1:1"
	valueProps := "font1:9:000:left:1:1:1:1"
	headerProps := "font1:10:100:left:1:1:1:1"
	radioProps := "font1:9:000:center:1:1:1:1"
	checkboxProps := "font1:8:000:center:1:1:1:1"
	smallTextProps := "font1:8:000:left:1:1:1:1"

	// 1. Configuration
	config := pdf.NewConfigBuilder().
		WithPage(pdf.PageSizeA4).
		WithPageBorder(1, 1, 1, 1).
		WithPageAlignment(1).
		Build()

	docBuilder := pdf.NewDocumentBuilder().
		WithConfig(config).
		WithTitle("font1:16:100:left:0:0:0:1", "   PATIENT REGISTRATION FORM").
		WithTitleTable(pdf.NewTableBuilder().
			WithColumns(3, []float64{0.3333333333333333, 0.3333333333333333, 0.3333333333333333}).
			AddRow(
				pdf.NewCell("font1:12:000:left:0:0:0:0", ""),
				pdf.NewCell("font1:18:100:center:0:0:0:0", "   PATIENT REGISTRATION FORM"),
				pdf.NewCell("font1:12:000:left:0:0:0:0", ""),
			).
			Build())

	// 2. Section A: Patient Information
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION A: PATIENT INFORMATION")).
		Build())

	// Name fields
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1, 2, 1, 2}).
		AddRow(
			pdf.NewCell(labelProps, "First Name:"),
			pdf.NewTextFieldCell(valueProps, "Michael", "first_name", "Michael"),
			pdf.NewCell(labelProps, "Last Name:"),
			pdf.NewTextFieldCell(valueProps, "Thompson", "last_name", "Thompson"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Middle Name:"),
			pdf.NewTextFieldCell(valueProps, "James", "middle_name", "James"),
			pdf.NewCell(labelProps, "Date of Birth:"),
			pdf.NewTextFieldCell(valueProps, "03/15/1985", "dob", "03/15/1985"),
		).
		AddRow(
			pdf.NewCell(labelProps, "SSN:"),
			pdf.NewTextFieldCell(valueProps, "***-**-4589", "ssn", "***-**-4589"),
			pdf.NewCell(labelProps, "Patient ID:"),
			pdf.NewTextFieldCell(valueProps, "PT-2024-78542", "patient_id", "PT-2024-78542"),
		).
		Build())

	// Gender (Radio buttons)
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(6, []float64{1, 0.5, 1, 0.5, 1, 0.5}).
		AddRow(
			pdf.NewCell(labelProps, "Gender:"),
			pdf.NewRadioCell(radioProps, "gender_male", "male", "gender", true),
			pdf.NewCell(valueProps, "Male"),
			pdf.NewRadioCell(radioProps, "gender_female", "female", "gender", false),
			pdf.NewCell(valueProps, "Female"),
			pdf.NewRadioCell(radioProps, "gender_other", "other", "gender", false),
		).
		Build())

	// Marital Status
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(8, []float64{1.5, 0.5, 1, 0.5, 1, 0.5, 1, 0.5}).
		AddRow(
			pdf.NewCell(labelProps, "Marital Status:"),
			pdf.NewRadioCell(radioProps, "marital_single", "single", "marital_status", false),
			pdf.NewCell(valueProps, "Single"),
			pdf.NewRadioCell(radioProps, "marital_married", "married", "marital_status", true),
			pdf.NewCell(valueProps, "Married"),
			pdf.NewRadioCell(radioProps, "marital_divorced", "divorced", "marital_status", false),
			pdf.NewCell(valueProps, "Divorced"),
			pdf.NewRadioCell(radioProps, "marital_widowed", "widowed", "marital_status", false),
		).
		Build())

	// 3. Section B: Contact Information
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION B: CONTACT INFORMATION")).
		Build())

	// Address
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(2, []float64{1, 3}).
		AddRow(
			pdf.NewCell(labelProps, "Street Address:"),
			pdf.NewTextFieldCell(valueProps, "4521 Oak Ridge Boulevard, Apt 12B", "street_address", "4521 Oak Ridge Boulevard, Apt 12B"),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(2, []float64{1, 3}).
		AddRow(
			pdf.NewCell(labelProps, "City:"),
			pdf.NewTextFieldCell(valueProps, "Austin", "city", "Austin"),
		).
		AddRow(
			pdf.NewCell(labelProps, "State:"),
			pdf.NewTextFieldCell(valueProps, "Texas", "state", "Texas"),
		).
		AddRow(
			pdf.NewCell(labelProps, "ZIP Code:"),
			pdf.NewTextFieldCell(valueProps, "78745", "zip_code", "78745"),
		).
		Build())

	// Phones and Email
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1, 2, 1, 2}).
		AddRow(
			pdf.NewCell(labelProps, "Home Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-1234", "home_phone", "(512) 555-1234"),
			pdf.NewCell(labelProps, "Cell Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-9876", "cell_phone", "(512) 555-9876"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Work Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-4567", "work_phone", "(512) 555-4567"),
			pdf.NewCell(labelProps, "Email:"),
			pdf.NewTextFieldCell(valueProps, "m.thompson@email.com", "email", "m.thompson@email.com"),
		).
		Build())

	// Section C: Emergency Contact
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION C: EMERGENCY CONTACT")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1, 2, 1, 2}).
		AddRow(
			pdf.NewCell(labelProps, "Contact Name:"),
			pdf.NewTextFieldCell(valueProps, "Sarah Thompson", "emergency_name", "Sarah Thompson"),
			pdf.NewCell(labelProps, "Relationship:"),
			pdf.NewTextFieldCell(valueProps, "Spouse", "emergency_relationship", "Spouse"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Phone Number:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-2468", "emergency_phone", "(512) 555-2468"),
			pdf.NewCell(labelProps, "Alt. Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-1357", "emergency_alt_phone", "(512) 555-1357"),
		).
		Build())

	// Section D: Insurance Information
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRowWithHeight(25, pdf.NewCell("font1:10:100:left:1:1:1:0", "SECTION D: INSURANCE INFORMATION")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1.5, 2, 1.5, 2}).
		AddRow(
			pdf.NewCell(labelProps, "Insurance Company:"),
			pdf.NewTextFieldCell(valueProps, "Blue Cross Blue Shield", "insurance_company", "Blue Cross Blue Shield"),
			pdf.NewCell(labelProps, "Policy Number:"),
			pdf.NewTextFieldCell(valueProps, "BCB-78542136", "policy_number", "BCB-78542136"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Group Number:"),
			pdf.NewTextFieldCell(valueProps, "GRP-45892", "group_number", "GRP-45892"),
			pdf.NewCell(labelProps, "Subscriber ID:"),
			pdf.NewTextFieldCell(valueProps, "SUB-123456789", "subscriber_id", "SUB-123456789"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Subscriber Name:"),
			pdf.NewTextFieldCell(valueProps, "Michael J. Thompson", "subscriber_name", "Michael J. Thompson"),
			pdf.NewCell(labelProps, "Subscriber DOB:"),
			pdf.NewTextFieldCell(valueProps, "03/15/1985", "subscriber_dob", "03/15/1985"),
		).
		Build())

	// 4. Section E: Medical History (Checkboxes)
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION E: MEDICAL HISTORY")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1.5, 2, 1.5, 2}).
		AddRow(
			pdf.NewCell(labelProps, "Primary Physician:"),
			pdf.NewTextFieldCell(valueProps, "Dr. Robert Williams", "primary_physician", "Dr. Robert Williams"),
			pdf.NewCell(labelProps, "Physician Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-8900", "physician_phone", "(512) 555-8900"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Preferred Pharmacy:"),
			pdf.NewTextFieldCell(valueProps, "CVS Pharmacy #4521", "preferred_pharmacy", "CVS Pharmacy #4521"),
			pdf.NewCell(labelProps, "Pharmacy Phone:"),
			pdf.NewTextFieldCell(valueProps, "(512) 555-7890", "pharmacy_phone", "(512) 555-7890"),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(labelProps, "Do you have any of the following conditions? (Check all that apply)")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(6, []float64{0.3, 1.2, 0.3, 1.2, 0.3, 1.2}).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "condition_diabetes", "diabetes", true),
			pdf.NewCell(smallTextProps, "Diabetes"),
			pdf.NewCheckboxCell(checkboxProps, "condition_hypertension", "hypertension", true),
			pdf.NewCell(smallTextProps, "Hypertension"),
			pdf.NewCheckboxCell(checkboxProps, "condition_heart_disease", "heart_disease", false),
			pdf.NewCell(smallTextProps, "Heart Disease"),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "condition_asthma", "asthma", false),
			pdf.NewCell(smallTextProps, "Asthma"),
			pdf.NewCheckboxCell(checkboxProps, "condition_arthritis", "arthritis", false),
			pdf.NewCell(smallTextProps, "Arthritis"),
			pdf.NewCheckboxCell(checkboxProps, "condition_cancer", "cancer", false),
			pdf.NewCell(smallTextProps, "Cancer"),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "condition_thyroid", "thyroid", false),
			pdf.NewCell(smallTextProps, "Thyroid Disorder"),
			pdf.NewCheckboxCell(checkboxProps, "condition_kidney", "kidney", false),
			pdf.NewCell(smallTextProps, "Kidney Disease"),
			pdf.NewCheckboxCell(checkboxProps, "condition_liver", "liver", false),
			pdf.NewCell(smallTextProps, "Liver Disease"),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "condition_depression", "depression", false),
			pdf.NewCell(smallTextProps, "Depression"),
			pdf.NewCheckboxCell(checkboxProps, "condition_anxiety", "anxiety", false),
			pdf.NewCell(smallTextProps, "Anxiety"),
			pdf.NewCheckboxCell(checkboxProps, "condition_other", "other", false),
			pdf.NewCell(smallTextProps, "Other"),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(2, []float64{1, 4}).
		AddRow(
			pdf.NewCell(labelProps, "Current Medications:"),
			pdf.NewTextFieldCell(valueProps, "Metformin 500mg, Lisinopril 10mg, Aspirin 81mg", "current_medications", "Metformin 500mg, Lisinopril 10mg, Aspirin 81mg"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Known Allergies:"),
			pdf.NewTextFieldCell(valueProps, "Penicillin, Shellfish", "known_allergies", "Penicillin, Shellfish"),
		).
		Build())

	// 5. Section F: Lifestyle (Radios)
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION F: LIFESTYLE INFORMATION")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(6, []float64{1.5, 0.5, 1, 0.5, 1, 0.5}).
		AddRow(
			pdf.NewCell(labelProps, "Do you smoke?"),
			pdf.NewRadioCell(radioProps, "smoke_yes", "yes", "smoking", false),
			pdf.NewCell(valueProps, "Yes"),
			pdf.NewRadioCell(radioProps, "smoke_no", "no", "smoking", true),
			pdf.NewCell(valueProps, "No"),
			pdf.NewRadioCell(radioProps, "smoke_former", "former", "smoking", false),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(6, []float64{1.5, 0.5, 1, 0.5, 1, 0.5}).
		AddRow(
			pdf.NewCell(labelProps, "Do you drink alcohol?"),
			pdf.NewRadioCell(radioProps, "alcohol_yes", "yes", "alcohol", false),
			pdf.NewCell(valueProps, "Yes"),
			pdf.NewRadioCell(radioProps, "alcohol_no", "no", "alcohol", false),
			pdf.NewCell(valueProps, "No"),
			pdf.NewRadioCell(radioProps, "alcohol_occasional", "occasional", "alcohol", true),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(6, []float64{1.5, 0.5, 1, 0.5, 1, 0.5}).
		AddRow(
			pdf.NewCell(labelProps, "Do you exercise?"),
			pdf.NewRadioCell(radioProps, "exercise_regular", "regular", "exercise", true),
			pdf.NewCell(valueProps, "Regularly"),
			pdf.NewRadioCell(radioProps, "exercise_occasional", "occasional", "exercise", false),
			pdf.NewCell(valueProps, "Occasionally"),
			pdf.NewRadioCell(radioProps, "exercise_never", "never", "exercise", false),
		).
		Build())

	// 6. Section G: Reason for Visit
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION G: REASON FOR VISIT")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(2, []float64{1.5, 4}).
		AddRow(
			pdf.NewCell("font1:9:100:left:1:1:1:0", "Reason for Visit:"),
			pdf.NewTextFieldCell("font1:9:000:left:1:1:1:0", "Annual Physical Examination", "reason_for_visit", "Annual Physical Examination"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Symptoms Description:"),
			pdf.NewTextFieldCell(valueProps, "Occasional fatigue, routine checkup for diabetes management", "symptoms_description", "Occasional fatigue, routine checkup for diabetes management"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Symptom Duration:"),
			pdf.NewTextFieldCell(valueProps, "2-3 weeks", "symptom_duration", "2-3 weeks"),
		).
		Build())

	// 7. Section H: Consent & Authorization
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell(headerProps, "SECTION H: CONSENT & AUTHORIZATION")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(2, []float64{0.2, 5}).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "consent_treatment", "consent_treatment", true),
			pdf.NewCell(smallTextProps, "I consent to receive medical treatment as deemed necessary by my healthcare provider."),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "consent_hipaa", "consent_hipaa", true),
			pdf.NewCell(smallTextProps, "I acknowledge receipt of the Notice of Privacy Practices (HIPAA)."),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "consent_billing", "consent_billing", true),
			pdf.NewCell(smallTextProps, "I authorize the release of medical information for billing and insurance purposes."),
		).
		AddRow(
			pdf.NewCheckboxCell(checkboxProps, "consent_financial", "consent_financial", true),
			pdf.NewCell(smallTextProps, "I accept financial responsibility for charges not covered by my insurance."),
		).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1.2, 2.5, 0.8, 1.5}).
		AddRow(
			pdf.NewCell(labelProps, "Patient Signature:"),
			pdf.NewTextFieldCell(valueProps, "Michael J. Thompson", "patient_signature", "Michael J. Thompson"),
			pdf.NewCell(labelProps, "Date:"),
			pdf.NewTextFieldCell(valueProps, "11/26/2025", "signature_date", "11/26/2025"),
		).
		AddRow(
			pdf.NewCell(labelProps, "Guardian Name:"),
			pdf.NewTextFieldCell(valueProps, "N/A", "guardian_name", "N/A"),
			pdf.NewCell(labelProps, "Relationship:"),
			pdf.NewTextFieldCell(valueProps, "N/A", "guardian_relationship", "N/A"),
		).
		Build())

	// 8. Office Use Only
	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(1, []float64{1}).
		AddRow(pdf.NewCell("font1:7:010:left:1:1:1:1", "FOR OFFICE USE ONLY")).
		Build())

	docBuilder.AddTable(pdf.NewTableBuilder().
		WithColumns(4, []float64{1, 2, 1, 2}).
		AddRow(
			pdf.NewCell("font1:8:100:left:1:1:1:1", "Received By:"),
			pdf.NewTextFieldCell("font1:8:000:left:1:1:1:1", "Jane Smith, RN", "received_by", "Jane Smith, RN"),
			pdf.NewCell("font1:8:100:left:1:1:1:1", "Date/Time:"),
			pdf.NewTextFieldCell("font1:8:000:left:1:1:1:1", "11/26/2025 09:30 AM", "received_datetime", "11/26/2025 09:30 AM"),
		).
		AddRow(
			pdf.NewCell("font1:8:100:left:1:1:1:1", "Verified By:"),
			pdf.NewTextFieldCell("font1:8:000:left:1:1:1:1", "Mary Johnson", "verified_by", "Mary Johnson"),
			pdf.NewCell("font1:8:100:left:1:1:1:1", "MRN:"),
			pdf.NewTextFieldCell("font1:8:000:left:1:1:1:1", "MRN-2024-785421", "mrn", "MRN-2024-785421"),
		).
		Build())

	// Footer
	docBuilder.WithFooter("font1:7:000:center", "CONFIDENTIAL PATIENT INFORMATION - PROTECTED UNDER HIPAA")

	return docBuilder.Build()
}
