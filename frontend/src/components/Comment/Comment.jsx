import image from "/src/static/img/x-button.png";
import "./Comment.css"
const Comment= ({comment, onDelete}) => {

    const handleDeleteCommentClick = () => {
        onDelete(comment._id);
    };

    return (
        <div className="comment-info" >
            <div className="comment-content">
                <div className="comment-user">
                    <img className="comment-user-image" src={comment.User.image} alt="image" />
                    <p>{comment.User.username}</p>
                </div>
                <img className="delete-comment" onClick={handleDeleteCommentClick} src={image} alt="delete" />
            </div>
            <div className="create-comment">
                <p>{comment.message}</p>
            </div>
        </div>
    );
};

export default Comment;