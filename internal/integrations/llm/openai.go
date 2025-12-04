package llm

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

type OpenAIProvider struct {
	ApiKey string
	Model  string
	Client *openai.Client
}

func NewOpenAIProvider(apiKey string, model string) *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAIProvider{
		ApiKey: apiKey,
		Model:  model,
		Client: &client,
	}
}

/* type ResponseNewParams struct {
    // Whether to run the model response in the background.
    // [Learn more](https://platform.openai.com/docs/guides/background).
    Background param.Opt[bool] `json:"background,omitzero"`
    // A system (or developer) message inserted into the model's context.
    //
    // When using along with `previous_response_id`, the instructions from a previous
    // response will not be carried over to the next response. This makes it simple to
    // swap out system (or developer) messages in new responses.
    Instructions param.Opt[string] `json:"instructions,omitzero"`
    // An upper bound for the number of tokens that can be generated for a response,
    // including visible output tokens and
    // [reasoning tokens](https://platform.openai.com/docs/guides/reasoning).
    MaxOutputTokens param.Opt[int64] `json:"max_output_tokens,omitzero"`
    // The maximum number of total calls to built-in tools that can be processed in a
    // response. This maximum number applies across all built-in tool calls, not per
    // individual tool. Any further attempts to call a tool by the model will be
    // ignored.
    MaxToolCalls param.Opt[int64] `json:"max_tool_calls,omitzero"`
    // Whether to allow the model to run tool calls in parallel.
    ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
    // The unique ID of the previous response to the model. Use this to create
    // multi-turn conversations. Learn more about
    // [conversation state](https://platform.openai.com/docs/guides/conversation-state).
    PreviousResponseID param.Opt[string] `json:"previous_response_id,omitzero"`
    // Whether to store the generated model response for later retrieval via API.
    Store param.Opt[bool] `json:"store,omitzero"`
    // What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
    // make the output more random, while lower values like 0.2 will make it more
    // focused and deterministic. We generally recommend altering this or `top_p` but
    // not both.
    Temperature param.Opt[float64] `json:"temperature,omitzero"`
    // An integer between 0 and 20 specifying the number of most likely tokens to
    // return at each token position, each with an associated log probability.
    TopLogprobs param.Opt[int64] `json:"top_logprobs,omitzero"`
    // An alternative to sampling with temperature, called nucleus sampling, where the
    // model considers the results of the tokens with top_p probability mass. So 0.1
    // means only the tokens comprising the top 10% probability mass are considered.
    //
    // We generally recommend altering this or `temperature` but not both.
    TopP param.Opt[float64] `json:"top_p,omitzero"`
    // Used by OpenAI to cache responses for similar requests to optimize your cache
    // hit rates. Replaces the `user` field.
    // [Learn more](https://platform.openai.com/docs/guides/prompt-caching).
    PromptCacheKey param.Opt[string] `json:"prompt_cache_key,omitzero"`
    // A stable identifier used to help detect users of your application that may be
    // violating OpenAI's usage policies. The IDs should be a string that uniquely
    // identifies each user. We recommend hashing their username or email address, in
    // order to avoid sending us any identifying information.
    // [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
    SafetyIdentifier param.Opt[string] `json:"safety_identifier,omitzero"`
    // This field is being replaced by `safety_identifier` and `prompt_cache_key`. Use
    // `prompt_cache_key` instead to maintain caching optimizations. A stable
    // identifier for your end-users. Used to boost cache hit rates by better bucketing
    // similar requests and to help OpenAI detect and prevent abuse.
    // [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#safety-identifiers).
    User param.Opt[string] `json:"user,omitzero"`
    // Specify additional output data to include in the model response. Currently
    // supported values are:
    //
    //   - `code_interpreter_call.outputs`: Includes the outputs of python code execution
    //     in code interpreter tool call items.
    //   - `computer_call_output.output.image_url`: Include image urls from the computer
    //     call output.
    //   - `file_search_call.results`: Include the search results of the file search tool
    //     call.
    //   - `message.input_image.image_url`: Include image urls from the input message.
    //   - `message.output_text.logprobs`: Include logprobs with assistant messages.
    //   - `reasoning.encrypted_content`: Includes an encrypted version of reasoning
    //     tokens in reasoning item outputs. This enables reasoning items to be used in
    //     multi-turn conversations when using the Responses API statelessly (like when
    //     the `store` parameter is set to `false`, or when an organization is enrolled
    //     in the zero data retention program).
    Include []ResponseIncludable `json:"include,omitzero"`
    // Set of 16 key-value pairs that can be attached to an object. This can be useful
    // for storing additional information about the object in a structured format, and
    // querying for objects via API or the dashboard.
    //
    // Keys are strings with a maximum length of 64 characters. Values are strings with
    // a maximum length of 512 characters.
    Metadata shared.Metadata `json:"metadata,omitzero"`
    // Reference to a prompt template and its variables.
    // [Learn more](https://platform.openai.com/docs/guides/text?api-mode=responses#reusable-prompts).
    Prompt ResponsePromptParam `json:"prompt,omitzero"`
    // Specifies the processing type used for serving the request.
    //
    //   - If set to 'auto', then the request will be processed with the service tier
    //     configured in the Project settings. Unless otherwise configured, the Project
    //     will use 'default'.
    //   - If set to 'default', then the request will be processed with the standard
    //     pricing and performance for the selected model.
    //   - If set to '[flex](https://platform.openai.com/docs/guides/flex-processing)' or
    //     'priority', then the request will be processed with the corresponding service
    //     tier. [Contact sales](https://openai.com/contact-sales) to learn more about
    //     Priority processing.
    //   - When not set, the default behavior is 'auto'.
    //
    // When the `service_tier` parameter is set, the response body will include the
    // `service_tier` value based on the processing mode actually used to serve the
    // request. This response value may be different from the value set in the
    // parameter.
    //
    // Any of "auto", "default", "flex", "scale", "priority".
    ServiceTier ResponseNewParamsServiceTier `json:"service_tier,omitzero"`
    // The truncation strategy to use for the model response.
    //
    //   - `auto`: If the context of this response and previous ones exceeds the model's
    //     context window size, the model will truncate the response to fit the context
    //     window by dropping input items in the middle of the conversation.
    //   - `disabled` (default): If a model response will exceed the context window size
    //     for a model, the request will fail with a 400 error.
    //
    // Any of "auto", "disabled".
    Truncation ResponseNewParamsTruncation `json:"truncation,omitzero"`
    // Text, image, or file inputs to the model, used to generate a response.
    //
    // Learn more:
    //
    // - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
    // - [Image inputs](https://platform.openai.com/docs/guides/images)
    // - [File inputs](https://platform.openai.com/docs/guides/pdf-files)
    // - [Conversation state](https://platform.openai.com/docs/guides/conversation-state)
    // - [Function calling](https://platform.openai.com/docs/guides/function-calling)
    Input ResponseNewParamsInputUnion `json:"input,omitzero"`
    // Model ID used to generate the response, like `gpt-4o` or `o3`. OpenAI offers a
    // wide range of models with different capabilities, performance characteristics,
    // and price points. Refer to the
    // [model guide](https://platform.openai.com/docs/models) to browse and compare
    // available models.
    Model shared.ResponsesModel `json:"model,omitzero"`
    // **o-series models only**
    //
    // Configuration options for
    // [reasoning models](https://platform.openai.com/docs/guides/reasoning).
    Reasoning shared.ReasoningParam `json:"reasoning,omitzero"`
    // Configuration options for a text response from the model. Can be plain text or
    // structured JSON data. Learn more:
    //
    // - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
    // - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
    Text ResponseTextConfigParam `json:"text,omitzero"`
    // How the model should select which tool (or tools) to use when generating a
    // response. See the `tools` parameter to see how to specify which tools the model
    // can call.
    ToolChoice ResponseNewParamsToolChoiceUnion `json:"tool_choice,omitzero"`
    // An array of tools the model may call while generating a response. You can
    // specify which tool to use by setting the `tool_choice` parameter.
    //
    // The two categories of tools you can provide the model are:
    //
    //   - **Built-in tools**: Tools that are provided by OpenAI that extend the model's
    //     capabilities, like
    //     [web search](https://platform.openai.com/docs/guides/tools-web-search) or
    //     [file search](https://platform.openai.com/docs/guides/tools-file-search).
    //     Learn more about
    //     [built-in tools](https://platform.openai.com/docs/guides/tools).
    //   - **Function calls (custom tools)**: Functions that are defined by you, enabling
    //     the model to call your own code. Learn more about
    //     [function calling](https://platform.openai.com/docs/guides/function-calling).
    Tools []ToolUnionParam `json:"tools,omitzero"`
    paramObj
} */

func (p *OpenAIProvider) GenerateResponse(prompt string, previousID string) (string, string, error) {
	params := responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model: p.Model,
		Store: openai.Bool(true),
	}

	if len(previousID) > 0 {
		params.PreviousResponseID = openai.String(previousID)
	}

	resp, err := p.Client.Responses.New(context.TODO(), params)

	if err != nil {
		panic(err.Error())
	}

	return resp.OutputText(), resp.ID, nil
}

func (p *OpenAIProvider) GenerateEmbedding(text string) ([]float32, error) {
	// TODO: Implement OpenAI embedding API
	return nil, fmt.Errorf("openai embedding not implemented")
}
